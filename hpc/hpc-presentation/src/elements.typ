#import "colors.typ": *

#let unbreak(body) = {
    set text(hyphenate: false)
    body
}

#let important-box(body,
    title: none,
    time: none,
    primary-color: black,
    secondary-color: black.lighten(90%),
    tertiary-color: black,
    dotted: false,
    figured: false,
    counter: none,
    show-numbering: false,
    numbering-format: (..n) => numbering("1.1", ..n),
    figure-supplement: none,
    figure-kind: none) = {
        let body = {
            if show-numbering or figured {
                counter.step()
            }
            set par(justify: true)
            show: align.with(left)
            block(width: 100%,
                inset: 1em,
                fill: secondary-color,
                stroke: (
                    thickness: 5pt,
                    paint: primary-color,
                    dash: if dotted { "dotted" } else { "solid" }
                ),
                text(size: 0.75em, strong(text(fill: tertiary-color, smallcaps(title) + if show-numbering or figured {
                    [ ] + context numbering(numbering-format, ..counter.at(here()))
                } + h(1fr) + time)))
                + block(body)
            )
        }
    body
}

#let standard-box-translations = (
    "definition": [Definizione],
    "example": [Esempio],
    "task": [Obbiettivo],
    "solution": [Soluzione],
    "hint": [Osservazione],
    "consequence": [Conseguenza],
    "notice": [Attenzione],
    "experiment": [Dati Sperimentali],
    "code": [Codice],
    "conclusion": [Conclusione],
)

#let definition = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("definition"),
    primary-color: green,
    secondary-color: green.lighten(90%),
    tertiary-color: green,
    figure-kind: "definition",
    counter: counter("grape-suite-element-definition"))

#let example = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("example"),
    primary-color: green,
    secondary-color: green.lighten(90%),
    tertiary-color: green,
    dotted: true,
    figure-kind: "example",
    counter: counter("grape-suite-element-example"))

#let task = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("task"),
    primary-color: blue,
    secondary-color: blue.lighten(90%),
    tertiary-color: blue,
    figure-kind: "task",
    counter: counter("grape-suite-element-task"))

#let solution = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("solution"),
    primary-color: darkblue,
    secondary-color: darkblue.lighten(90%),
    tertiary-color: darkblue,
    figure-kind: "solution",
    counter: counter("grape-suite-element-solution"))

#let hint = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("hint"),
    primary-color: orange,
    secondary-color: orange.lighten(90%),
    tertiary-color: orange,
    figure-kind: "hint",
    counter: counter("grape-suite-element-hint"))

#let consequnce = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("consequnce"),
    primary-color: brown,
    secondary-color: brown.lighten(90%),
    tertiary-color: brown,
    figure-kind: "consequnce",
    counter: counter("grape-suite-element-consequnce"))

#let notice = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("notice"),
    primary-color: red,
    secondary-color: red.lighten(90%),
    tertiary-color: red,
    figure-kind: "notice",
    counter: counter("grape-suite-element-notice"))

#let experiment = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("experiment"),
    primary-color: green,
    secondary-color: green.lighten(90%),
    tertiary-color: green,
    dotted: true,
    figure-kind: "experiment",
    counter: counter("grape-suite-element-experiment"))

#let code = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("code"),
    primary-color: black,
    secondary-color: black.lighten(90%),
    tertiary-color: black,
    figure-kind: "code",
    counter: counter("grape-suite-element-code"))

#let conclusion = important-box.with(
    title: context state("grape-suite-box-translations", standard-box-translations).final().at("conclusion"),
    primary-color: brown,
    secondary-color: brown.lighten(90%),
    tertiary-color: brown,
    figure-kind: "conclusion",
    counter: counter("grape-suite-element-conclusion"))

#let sentence-logic(body) = {
    show figure.where(kind: "example"): it => {
        show: pad.with(0.25em)
        grid(columns: (1cm, 1fr),
            column-gutter: 0.5em,
            context [(#(counter("grape-suite-sentence-counter").at(here()).first()+1))],
            it.body)
    }
    body
}

#let sentence(body) = {
    figure(kind: "example", supplement: context state("grape-suite-element-sentence-supplement", "Example").final(), align(left, body) +
    counter("grape-suite-sentence-counter").step())
}

#let blockquote(body, source) = pad(x: 1em, y: 0.25em, body + block(text(size: 0.75em, source)))