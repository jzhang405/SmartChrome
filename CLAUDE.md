# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This repository uses the **Specify framework** - a structured workflow system for feature development. It's currently configured as a template project with no actual source code yet, but contains the complete Specify workflow infrastructure.

## Key Architecture

### Specify Framework Structure

- **`.specify/`**: Core framework templates and scripts
  - `templates/`: Document templates for spec, plan, tasks, etc.
  - `scripts/bash/`: Automation scripts for workflow steps
  - `memory/constitution.md`: Project principles and constraints (currently template)

- **`.claude/commands/`**: Claude Code slash commands for the Specify workflow
  - `/speckit.specify`: Create feature specifications
  - `/speckit.plan`: Generate implementation plans
  - `/speckit.tasks`: Create task breakdowns
  - `/speckit.implement`: Execute implementation
  - `/speckit.clarify`: Resolve specification ambiguities
  - `/speckit.analyze`: Cross-artifact consistency analysis
  - `/speckit.constitution`: Manage project constitution
  - `/speckit.checklist`: Generate validation checklists

### Workflow Process

1. **Specification** (`/speckit.specify`): Create feature specs in `specs/[###-feature-name]/spec.md`
2. **Planning** (`/speckit.plan`): Generate technical plans with research, data models, contracts
3. **Task Generation** (`/speckit.tasks`): Create implementation task breakdown
4. **Implementation** (`/speckit.implement`): Execute tasks following TDD approach

## Development Commands

### Feature Development Workflow

```bash
# Start a new feature
/speckit.specify "Feature description here"

# Generate implementation plan
/speckit.plan

# Create task breakdown
/speckit.tasks

# Execute implementation
/speckit.implement
```

### Project Setup

When starting development on this project:
1. First define the project constitution using `/speckit.constitution`
2. Use the slash commands above for all feature development
3. The framework will guide you through the complete development lifecycle

## Important Notes

- This is currently a **template repository** - no source code exists yet
- All development should follow the Specify workflow using the provided slash commands
- The constitution file (`.specify/memory/constitution.md`) should be populated with project-specific principles before development begins
- Feature specifications, plans, and tasks are stored in the `specs/` directory
- The framework enforces TDD (Test-Driven Development) and independent user story implementation
