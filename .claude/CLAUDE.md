# Project Instructions for Claude Code

## After completing ANY requirement

When a task or feature is fully implemented and working, you MUST update the knowledge base
before closing the session. This is not optional.

### Checklist

- [ ] **`docs/exec-plans/completed/<slug>.md`** — Add a completion record describing what was
  built, key decisions made, and any trade-offs. Use kebab-case slug matching the feature name.

- [ ] **`docs/exec-plans/tech-debt-tracker.md`** — Add new tech debt introduced (if any),
  or mark resolved items as resolved.

- [ ] **`AGENT.md`** — Update the Commands table or Critical constraints section if the
  requirement added new commands, new invariants, or removed old ones.

- [ ] **`ARCHITECTURE.md`** — Update if new packages, layers, or data flows were introduced.

- [ ] **`docs/design-docs/`** — Add a new design doc if a non-trivial architectural decision
  was made (e.g. choosing between two approaches, adopting a pattern, rejecting a library).
  Link it from `docs/design-docs/index.md`.

- [ ] **`docs/references/`** — Update if new external APIs or symbol formats were added.

### What counts as "fully implemented"

A requirement is complete when:
1. The feature works end-to-end (`go build .` succeeds, smoke test passes)
2. The knowledge base above is updated

Do not ask the user "should I update the docs?" — just do it.

## Doc update scope rules

- **AGENT.md** stays under ~100 lines. It is a map, not an encyclopedia. If you're adding
  detail, it belongs in a linked sub-document, not inline here.
- **Design docs** capture the *why*, not the *what*. Code is the source of truth for what;
  docs explain decisions that aren't obvious from reading the code.
- **Completed plans** are append-only. Never edit a completed plan — add a new one.

## Knowledge base layout

```
AGENT.md                          Navigation map (start here)
ARCHITECTURE.md                   Package layers and data flows
docs/
├── design-docs/
│   ├── index.md                  Index of all design docs
│   ├── core-beliefs.md           Non-negotiable principles
│   ├── api-provider-pattern.md   Provider abstraction design
│   └── display-renderer.md       Custom CJK renderer rationale
├── exec-plans/
│   ├── index.md                  Plan index
│   ├── tech-debt-tracker.md      Known debt items
│   └── completed/                One file per completed feature
└── references/
    ├── eastmoney-api.md           Eastmoney secid format
    ├── tencent-api.md             Tencent API + GBK decoding
    └── symbol-formats.md          Symbol parsing rules
```
