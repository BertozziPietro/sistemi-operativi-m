#import "@preview/polylux:0.3.1"

#import "colors.typ": *
#import "dates.typ": semester, weekday
#import "elements.typ": *

#let uncover = polylux.uncover
#let only = polylux.uncover
#let pause = polylux.pause
#let slide(..args) = {
    state("grape-suite-slides", ()).update(k => {
        k.push("normal")
        k
    })
    polylux.polylux-slide(..args)
}

#let slides(
    series: none,
    title: none,
    topics: (),

    head-replacement: none,
    title-replacement: none,
    footer: none,

    author: none,
    lecturer: none,
    email: none,

    page-numbering: (n, total) => {
        text(size: 0.75em, strong[#n.first()])
        text(size: 0.5em, [ \/ #total.first()])
    },

    show-title-slide: true,
    show-author: true,
    show-lecturer: true, 
    show-semester: true,
    show-date: true,
    show-outline: true,
    show-todolist: true,
    show-footer: true,
    show-page-numbers: true,

    box-definition-title: standard-box-translations.at("definition"),
    box-example-title: standard-box-translations.at("example"),
    box-task-title: standard-box-translations.at("task"),
    box-solution-title: standard-box-translations.at("solution"),
    box-hint-title: standard-box-translations.at("hint"),
    box-consequence-title: standard-box-translations.at("consequence"),
    box-notice-title: standard-box-translations.at("notice"),
    box-experiment-title: standard-box-translations.at("experiment"),
    box-code-title: standard-box-translations.at("code"),
    box-conclusion-title: standard-box-translations.at("conclusion"),
    sentence-supplement: "Example",

    date: datetime.today(),
    body
) = {
        set text(hyphenate: false)
        let left-footer = if footer != none {
            footer
        } else {
            text(size: 0.5em, (
                if show-semester [#semester(short: true, date)],
                [#series],
                title,
                if show-author { author },
                if show-lecturer { lecturer }
              ).filter(e => e != none).join[ -- ]
            )
        }
    
    show footnote.entry: set text(size: 0.5em)
    
    show heading: set text(fill: darkblue)
    
    set text(size: 24pt, font: "Consolas")
    set page(paper: "presentation-16-9",
        footer: {
            let fs = state("grape-suite-slides", ())
            context {
                set text(fill: black)
                left-footer
                h(1fr)
                if show-page-numbers {
                    page-numbering(counter(page).at(here()), counter(page).final())
                }
            }
        }
    )
    
    state("grape-suite-box-translations").update((
        "definition": box-definition-title,
        "example": box-example-title,
        "task": box-task-title,
        "solution": box-solution-title,
        "hint": box-hint-title,
        "cosequence": box-consequence-title,
        "notice": box-notice-title,
        "experiment": box-experiment-title,
        "code": box-code-title,
        "conclusion": box-conclusion-title,
    ))
    
    state("grape-suite-element-sentence-supplement").update(sentence-supplement)
    show: sentence-logic
    
    if show-title-slide {
        slide(align(horizon, [
            #block(inset: (left: 1cm, top: 3cm))[
                #if head-replacement == none [
                    #text(fill: darkblue, size: 2em, strong[#series ]) \
                ] else { head-replacement }
                #if title-replacement == none [
                    #text(fill: darkblue.lighten(25%), strong(title))
                ] else { title-replacement }
    
                #set text(size: 0.75em)
                #if show-author [#author #if email != none [--- #email ] \ ]
                #if show-lecturer [#lecturer #if email != none [--- #email ] \ ]
                #if show-semester [#semester(date) \ ]
                #if show-date [#weekday(date.weekday()), #date.display("[day].[month].[year]")]
            ]
        ]))
    }

    if show-outline {
        set page(fill: darkblue, footer: context if show-footer {
            set text(fill: if here().page() > 2 or not show-outline { black } else {  white })
            left-footer
        })
    }
    counter(page).update(1)
    set page(fill: white)
    body
}

