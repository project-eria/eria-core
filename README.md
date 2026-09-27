# eria-core

Surcouche applicative d'ERIA au-dessus de [`go-wot`](../go-wot/). Fournit le
bootstrap (`Init`/`Start`), la configuration YAML, les automations, le scheduler
cron, ainsi que l'exposition et la consommation simplifiées de Things WoT.

Voir [CLAUDE.md](../CLAUDE.md) pour la vue d'ensemble du projet.

## EriaAppController

Chaque application qui expose des Things locaux obtient automatiquement un Thing
de contrôle, le **`EriaAppController`**. Il est ajouté par
[`AddAppController`](appController.go) dès le premier appel à `eria.Producer(instance)`
([producer.go](producer.go)) — aucune configuration n'est requise.

Il permet de piloter l'application à distance (aujourd'hui : ajuster le niveau
de log à chaud, sans redémarrage).

### Identité du Thing

- **ID (`urn`)** : `eria:app:controler:<binaire>@<machine>[:<instance>]` —
  unique sur le réseau (c'est aussi le nom de l'instance mDNS) ; `<instance>`
  est l'argument passé à `eria.Producer(instance)`, omis s'il est vide. Le
  préfixe `eria:app:controler:` est stable : l'intégration Home Assistant
  reconnaît les contrôleurs grâce à lui.
- **Référence / préfixe d'URL** : `eria` → toutes les routes sont servies sous
  `/eria`.
- **Sécurité** : `nosec_sc` (aucune authentification).

> Le préfixe `/eria` est identique pour toutes les instances ; seul l'ID du
> Thing varie via `<instance>`.

### Endpoints

Les URLs utilisent le `host`/`port` configurés (`port` par défaut : `80`, voir
[config.go](config.go)).

| Méthode | Chemin              | Opération WoT                     | Description                                    |
| ------- | ------------------- | --------------------------------- | ---------------------------------------------- |
| `GET`   | `/eria`             | —                                 | Thing Description JSON du contrôleur           |
| `GET`   | `/eria/logLevel`    | `readproperty`                    | Lit le niveau de log courant                   |
| `PUT`   | `/eria/logLevel`    | `writeproperty`                   | Change le niveau de log                        |
| `WS`    | `/eria/logLevel`    | `observeproperty`/`unobserve…`    | Observe les changements du niveau de log       |
| `GET`   | `/eria/automations` | `readproperty`                    | Liste JSON des automations planifiées          |

### Propriété `logLevel`

Niveau de log [zerolog](https://github.com/rs/zerolog), observable. Valeurs
autorisées (enum) :

```
error | warn | info | debug | trace
```

Une valeur hors de cette liste est rejetée par la validation du schéma go-wot.
L'écriture applique le niveau globalement et immédiatement
([appController.go](appController.go), `logLevelWrite`).

### Exemples

Lire le niveau courant :

```bash
curl http://localhost:80/eria/logLevel
# → "info"
```

Passer en `debug` (le corps est une chaîne JSON) :

```bash
curl -X PUT http://localhost:80/eria/logLevel \
     -H 'Content-Type: application/json' \
     -d '"debug"'
```

Récupérer la Thing Description complète :

```bash
curl http://localhost:80/eria
```

### Propriété `automations`

Propriété en lecture seule renvoyant la liste des automations actuellement
planifiées. La liste est vivante : elle change au gré des contextes actifs et
du recalcul des horaires (schedules `at|hour|…`).

Chaque entrée correspond à un groupe actif d'une automation :

| Champ        | Type              | Description                                                        |
| ------------ | ----------------- | ----------------------------------------------------------------- |
| `automation` | string            | Nom de l'automation (config YAML)                                 |
| `devices`    | []string          | Références des Things ciblés par l'action                         |
| `action`     | string            | Clé de l'action invoquée (ex. `on`, `toggle`)                     |
| `value`      | any (optionnel)   | Valeur de l'action, si présente                                   |
| `schedule`   | string            | Description lisible du schedule (ex. `every day at 07:30:00`)     |
| `nextRun`    | date-time (opt.)  | Prochaine exécution (horaires) ; absent pour `immediate`/désactivé |

```bash
curl http://localhost:80/eria/automations
```

```json
[
  {
    "automation": "Volets salon",
    "devices": ["shutter-salon"],
    "action": "close",
    "schedule": "every day at 21:30:00",
    "nextRun": "2026-07-01T21:30:00+02:00"
  }
]
```

## Génération des endpoints

Les routes ne sont pas déclarées à la main : elles sont dérivées de la Thing
Description par les bindings de `go-wot`.

- HTTP `GET`/`PUT` par propriété et `GET /<ref>` pour la TD :
  [go-wot/protocolHttp/httpServer.go](../go-wot/protocolHttp/httpServer.go).
- WebSocket `observe` par propriété (uniquement si la propriété est
  `Observable`, ce qui est le défaut) :
  [go-wot/protocolWebSocket/wsServer.go](../go-wot/protocolWebSocket/wsServer.go).

Ajouter une propriété ou une action au modèle `EriaAppController`
([model/models.go](model/models.go)) crée donc automatiquement les endpoints
correspondants.

## Découverte mDNS/DNS-SD

### Annonce

Toutes les apps annoncent leurs Things exposés en `_wot._tcp` (mDNS/DNS-SD,
conforme W3C WoT Discovery) dès `eria.Start()`. Configuration :

```yaml
mdns:
  disabled: false        # opt-out de l'annonce
  scheme: http           # http|https, annoncé dans le TXT (défaut http)
  discoveryTimeout: 60s  # délai de découverte des remoteThings `mdns:` requis
```

> Si `exposedAddr` (reverse proxy) est configuré, l'annonce mDNS pointera
> l'IP locale et non le proxy : mettre `mdns.disabled: true` dans ce cas.

Vérification : `dns-sd -B _wot._tcp local.` (macOS) ou
`avahi-browse -r _wot._tcp` (Linux).

### Connexion par identifiant mDNS

Une entrée `remoteThings` déclare **soit** `url:` **soit** `mdns:` :

```yaml
remoteThings:
- mdns: eria-aabbcc      # identifiant DNS-SD, XOR avec url
  tags: [ "temp" ]
  onConnect: required    # défaut — voir ci-dessous
  onLoss: rediscover     # défaut — voir ci-dessous
- url: http://manager-contexts.lan
  tags: [ "contexts" ]
```

L'identifiant `mdns` est comparé à l'**instance DNS-SD** (le TD id sans son
préfixe `urn:`, ex. `eria:test:1`), puis en repli au **hostname mDNS** du
device (label sans `.local.`, ex. `eria-aabbcc` pour un ESP32). L'URL du TD
est reconstruite depuis l'annonce (TXT `td`/`scheme` + adresse + port), puis
le pipeline de connexion habituel s'applique — le choix HTTP/WS par
interaction reste dicté par les Forms du TD.

Comportements par entrée (`url` comme `mdns`) :

- `onConnect: required` (défaut) — le démarrage attend la connexion ;
  introuvable après `mdns.discoveryTimeout` ⇒ Fatal.
- `onConnect: background` — l'app démarre sans attendre ; la connexion
  s'établit dès que le Thing apparaît (retry en tâche de fond).
- `onLoss: rediscover` (défaut) — quand une interaction échoue, re-résolution
  mDNS (l'IP a pu changer) puis reconnexion, en conservant le pointeur
  `EriaThing` détenu par l'app.
- `onLoss: ignore` — log Warn, reste déconnecté.
- `onLoss: fatal` — arrêt de l'app.

Limites actuelles : la perte est détectée passivement (erreur de
read/write/invoke — pas de heartbeat) ; les observations WebSocket établies
avant la perte continuent de viser l'ancienne adresse jusqu'au redémarrage.
