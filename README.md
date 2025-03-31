# Sjb site
Een nieuwe, snelle en veilige site voor Sint Jansbrug gemaakt met Go, Tailwind, Templ en HTMX.

* Go - Backend
* Tailwind - CSS
* Templ - Templating

## Makefile
De makefile zorgt ervoor dat je makkelijk de site kan draaien en bouwen. 

`make docker-dev` zorgt ervoor dat je niet elke keer opnieuw moet builden bij veranderingen.

Momenteel werkt alleen de dev omgeving, maar er is ook een productie omgeving.

## Ports

* 4000 - Site (80 in de dev env)
* 3333 - PhpMyAdmin (database shit)

## Dir structuur

`internal/`: Hier staat alle code buiten de router
`cmd/main.go`: Entrypoint van de site en waar de router staat
`static/`: static fileserver gaat hieruit
`dev/`: scriptjes en tools

## Veiligheid

Met gorm en templ zijn we volledig beschermd tegen sql injection en xss attacks.
Op alle niet-publieke routes staan middleware voor admin of leden privileges.

Useruploads worden "beschermd" via uuid, een hacker zou dus miljarden requests moeten doen om er een te vinden. (eigenlijk moet dit door middleware)
