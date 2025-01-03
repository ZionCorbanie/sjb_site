# Sjb site
Een nieuwe site voor Sint Jansbrug gemaakt met Go, Tailwind, Templ en HTMX.

* Go - Backend
* Tailwind - CSS
* Templ - Templating
* HTMX - Frontend zonder bullshit


## Makefile
De makefile zorgt ervoor dat je makkelijk de site kan draaien en bouwen. 

`make docker-dev` zorgt ervoor dat je niet elke keer opnieuw moet builden bij veranderingen.

## Ports

* 4000 - Site
* 3333 - PhpMyAdmin (database shit)

## Dir structuur

`internal/`: Hier staat alle code buiten de router
`cmd/main.go`: Entrypoint van de site en waar de router staat
`static/`: spreekt voor zichzelf, als je iets wil exposen doe hier
`dev/`: scriptjes en tools

## Naming

Omdat computer science kut is en ik mn leven haat is alles wat de user zou kunnen zien nederlands en anders engels. Endpoints zijn dus nederlands en data ook. Verder alles engels.

Verder idk camelCase/PascalCase omdat dat conanical go is. Pas op met json en sql want daar is de norm snake_case.

## Code style

Dus ja canon go pls. Werk zoveel mogelijk met de go features van een functie attachen aan een struct, zo is mocking heel nice voor als we tests zouden schrijven. Goed voorbeeld hiervan is store.go en de files in dbstore/ in static/store/

Verder in de handlers zie je voor elke endpoint een file, is winnen want anders kom je snel aan files met 1000+ lines.

Liefst geen css maar altijd tailwind. Ook alpine.js als er geen andere optie is via html en tailwind.

## Database

We gebruiken GORM voor de database. Die doet dus automatisch de structures uit store.go in de db. Maak snel ook ffe een mock db aan via phpmyadmin is niet zo moeilijk (die van mij is ook nogsteeds heel klein dus eh).
Je kan best 2 accounts registeren op de website en dan één admin maken via phpmyadmin.

We kunnen later wel een betere db en static/uploads dir maken maar voor nu is dit prima aangezien nog niet alle reqs goed zijn uitgewerkt.

## Hoe beginnen programmeren

* Leer go
* Leer htmx
* Leer tailwind (focus op mobile first design)
* Leer templ (ook niet moeilijk)
* Kijk naar de code en probeer te begrijpen (begin bij cmd/main.go en volg de router)
* Neem een van de (voor nu nog onbestaande) issues en maak branch
