---
name: "grill-with-docs"
description: "A relentless interview to sharpen a plan or design, which also creates docs (ADRs and glossary) as we go. Use when the user says grill with docs or wants a plan stress-tested and documented."
---
 
# Grill with docs
 
Source: mattpocock/skills (engineering/grill-with-docs). The original calls two skills, `grilling` and `domain-modeling`; both are bundled below so this skill is self-contained.
 
Run Part A (grilling) and Part B (domain modeling) at the same time: grill in rounds, and write the glossary and ADRs the moment terms and decisions crystallise.
 
## Part A — Grilling
 
Interview the user relentlessly until you reach a shared understanding. Map this as a **design tree**: every decision branches into the decisions that hang off it.
 
Work the tree in **rounds**. The **frontier** is every decision whose prerequisites are already settled: the questions you can ask _now_ without guessing at answers you haven't heard yet. Ask the whole frontier in one round: number each question and give your recommended answer. Then wait for the user's answers before the next round.
 
Format a round like so:
 
```
❓ **Q1** - **<question title>**: <question body, might be multiple paragraphs, including multiple choices>
 
➡️ <your recommended answer>
 
---
 
❓ **Q2** - **<question title>**: <question body>
 
➡️ <your recommended answer>
```
 
Word each question so "yes" accepts your recommended answer.
 
Each round the user answers reshapes the tree: settled decisions push the frontier outward and unblock questions that depended on them. Recompute the frontier and ask the next round. A question whose answer depends on another question still open in this round belongs to a _later_ round, not this one.
 
Finding _facts_ is your job, never the user's. When a frontier question needs a fact from the environment (filesystem, tools, web), dispatch a sub-agent with a fresh context window to find it; don't ask the user for anything you could look up yourself. Don't block on it: a running exploration is an unsettled prerequisite, so only the questions downstream of it wait; ask the rest of the frontier now. The _decisions_ are the user's: put each to them and wait.
 
The session is done when the frontier is empty: every branch of the design tree visited, nothing left silently assumed. Do not act on it until the user confirms you have reached a shared understanding.
 
## Part B — Domain modeling
 
Actively build and sharpen the project's domain model as you design: challenge terms, invent edge-case scenarios, and write the glossary and decisions down the moment they crystallise.
 
### File structure
 
Most projects have a single context: `GLOSSARY.md` at the root and ADRs in `docs/adr/0001-slug.md`, `0002-slug.md`… If a `GLOSSARY-MAP.md` exists at the root, there are multiple contexts, each with its own `GLOSSARY.md` and `docs/adr/`; the map lists the contexts, where they live and how they relate. Create files lazily, only when you have something to write.
 
### During the session
 
- **Challenge against the glossary.** When the user uses a term that conflicts with `GLOSSARY.md`, call it out immediately: "Your glossary defines X as A, but you seem to mean B. Which is it?"
- **Sharpen fuzzy language.** When a term is vague or overloaded, propose a precise canonical term.
- **Discuss concrete scenarios.** Invent edge-case scenarios that force precision about the boundaries between concepts.
- **Cross-reference with code.** When the user states how something works, check whether the code agrees and surface contradictions.
- **Update GLOSSARY.md inline** when a term is resolved; don't batch. The glossary holds no implementation details; it is a glossary and nothing else.
### GLOSSARY.md format
 
```md
# {Context Name}
 
{One or two sentences on what this context is and why it exists.}
 
## Language
 
**Order**:
{One or two sentence description}
_Avoid_: Purchase, transaction
```
 
Rules: be opinionated (pick the best word, list others under `_Avoid_`); one or two sentences, defining what it IS; only project-specific terms, never general programming concepts; group under subheadings when clusters emerge.
 
### ADRs — offer sparingly
 
Offer an ADR only when all three are true: (1) hard to reverse, (2) surprising without context, (3) the result of a real trade-off. Format:
 
```md
# {Short title of the decision}
 
{1-3 sentences: context, what we decided, and why.}
```
 
Optional sections only when they add value: Status (`proposed | accepted | deprecated | superseded by ADR-NNNN`), Considered Options, Consequences. Number by scanning `docs/adr/` for the highest number and incrementing.
 
What qualifies: architectural shape; integration patterns between contexts; technology choices with lock-in; boundary and scope decisions (explicit no-s count); deliberate deviations from the obvious path; constraints not visible in the code; non-obvious rejected alternatives.