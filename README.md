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

- **ID (`urn`)** : `eria:app:controler:<instance>` — `<instance>` est l'argument
  passé à `eria.Producer(instance)`.
- **Référence / préfixe d'URL** : `eria` → toutes les routes sont servies sous
  `/eria`.
- **Sécurité** : `nosec_sc` (aucune authentification).

> Le préfixe `/eria` est identique pour toutes les instances ; seul l'ID du
> Thing varie via `<instance>`.

### Endpoints

Les URLs utilisent le `host`/`port` configurés (`port` par défaut : `80`, voir
[config.go](config.go)).

| Méthode | Chemin           | Opération WoT                     | Description                                   |
| ------- | ---------------- | --------------------------------- | --------------------------------------------- |
| `GET`   | `/eria`          | —                                 | Thing Description JSON du contrôleur          |
| `GET`   | `/eria/logLevel` | `readproperty`                    | Lit le niveau de log courant                  |
| `PUT`   | `/eria/logLevel` | `writeproperty`                   | Change le niveau de log                       |
| `WS`    | `/eria/logLevel` | `observeproperty`/`unobserve…`    | Observe les changements du niveau de log      |

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
