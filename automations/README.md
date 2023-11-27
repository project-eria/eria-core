# Automation

- Note: Manual Override, Chaque automation peut-être manuellement désactivée (`none`), remplacée (nouvelle condition), pour 1 action, pour une période ou jusqu'a nouvelle modification.

## Conditions
- `context|<context name>`
	- not: `context|!<context name>`
- `time|[before|after]=<hour in 15:00 format>`
	?? Comment differencier "executer seulement si after..x" et "exectuer dès que after...x" ??

- `property|<param[==<>!=]value>` (opérateur de comparaison,  e.g. a==1)

## Actions
- `<action>|<param name>=<value>|<param name>=<value>`

## Scheduling
- `<empty>`: Immediate run
- `at|<time>`
	- fixed hour: `at|<hour in 15:00 format>`
	- thing hour property: `at|<thing device>=<thing property>`
		- min option: `at|<thing device>=<thing property>|min=<hour in 15:00 format>`
		- max option: `at|<thing device>=<thing property>|max=<hour in 15:00 format>`
??l'action est planifiée dans un cron, et une re-compilation tout les jours est nécessaire (par défaut matin à 3h, à cause du changement heure été|hivers, mais modifiable)
	- L'action est listée dans une liste 'scheduled'

- `every|<delay (10m, 1h, ...)>`, ??l'action est planifiée dans un cron, et une re-evaluation à chaque boucle
- `in|<delay (10m, 1h, ...)>`, ??l'action est planifiée et executée 1 fois

- ??Si `context` une re-evaluation à chaque changement
- ??Si `property`une re-evaluation à chaque changement

## e.g. Planned once, executed all times
- `name: <name description>`
- `action: <action id>|<param=value>|<param=value>`
- `condition:` e.g. `property|<device>=<property>|<value>`

```
- contexts: [away]
hour: none
- hour: 15:00
```

| Conditions                   | Scheduled  | Description |
| ---------------------------- | ---------- | ----------- |
| `context|!away`              | `at|15:00` |             |

Note: `context|!away` to disable the automation when context is `away`

### e.g. First condition 'true' is planned
Si plusieurs conditions groups, le premier `true` l'emporte
- `name: <name>` e.g. Ouverture Matin
- ?? `planning: <when to re-compile the condition>` e.g. every morning at 3am (cron format ?)
- `action: <action id>|<xxx>`
- ?? `disable:` e.g.  disable on certain context `context|away`
- `conditions list:` e.g. with shutters

```
- contexts: [away]
  hour: none
- contexts: [weekend, holiday]
  hour: goldenHourMorning
  min: 8:00
- hour: goldenHourMorning # monday, tuesday, wednesday, thursday, friday
  min: 7:15
```

Scheduled
| Conditions        | Scheduled                       | Description           |
| ----------------- | ------------------------------- | ----------------------|
| `context|weekend` | `at|astral=goldenHourMorning|min=8:00` | On WeekEnd            |
| `context|away`    | `at|astral=goldenHourMorning|min=8:00` | On Holiday            |
|                   | `at|astral=goldenHourMorning|min=7:15` | monday, ..., friday   |

Immediate
| Conditions        | Scheduled                       | Description           |
| ----------------- | ------------------------------- | ----------------------|
| `context|away`    | <empty>						  | On Holiday            |

?? Sub schedule conditions groups ??
