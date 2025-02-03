#let days = (
    "Lunedì",
    "Martedì",
    "Martedì",
    "Giovedì",
    "Venerdì",
    "Sabato",
    "Domenica"
)

#let weekday(short: false, daynr) = {
    let day = days.at(daynr - 1)
    if short {  day = day.slice(0, 2) }
    day
}

#let semester(short: false, date) = {
    let sem = (
        long: "Anno Accademico",
        short: "A.A."
    ).at(if short { "short" } else { "long" })

    let year = if date.month() < 4 { [#(date.year() - 1)/#date.year()] }
    else if date.month() >= 10 { [#date.year()/#(date.year() + 1)] }
    else if date.month() < 10 { date.year() }

    [#sem #year]
}