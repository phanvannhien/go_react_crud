# Agent Runtime Contract

## Execution Contract

Every request MUST pass through:

Phase 1: Intent Classification  
Phase 2: Environment Detection  
Phase 3: Risk Scoring  
Phase 3.5: Workflow Binding

Each valid entry point MUST bind to exactly one workflow:

create_feature
  → workflows/feature_workflow.md

modify_schema
  → workflows/migration_workflow.md

deploy_release
  → workflows/release_workflow.md

fix_production_bug
  → workflows/bugfix_workflow.md

optimize_performance
  → workflows/bugfix_workflow.md

investigate_security_issue
  → workflows/bugfix_workflow.md

create_ui_feature
  → workflows/ui_feature_workflow.md

modify_ui_component
  → workflows/ui_modification_workflow.md

fix_ui_bug
  → workflows/ui_bugfix_workflow.md

If no workflow is matched → Hard Stop.
Phase 4: Skill Routing  
Phase 5: Dependency Resolution  
Phase 6: Guard Enforcement  
Phase 7: Output Generation  
Phase 8: Post-Generation Validation  

If Phase 6 fails → Hard Stop  
If Phase 8 fails → Regenerate with corrections  

Risk Levels:

0–30   → Low  
31–60  → Moderate  
61–85  → High  
86–100 → Critical (Execution forbidden without override)