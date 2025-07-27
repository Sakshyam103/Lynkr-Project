# Git Practices Rule

## Branch Naming Conventions
- Feature branches: `feature/short-description`
- Bug fixes: `bugfix/issue-number-description`
- Hotfixes: `hotfix/issue-number-description`
- Release branches: `release/version-number`
- Documentation: `docs/what-changed`
- Refactoring: `refactor/what-changed`

## Commit Message Standards
- Use the imperative mood: "Add feature" not "Added feature"
- Start with a capital letter
- Keep the subject line under 50 characters
- Use the body to explain what and why, not how
- Reference issue numbers when applicable: "Fix login bug (fixes #123)"
- Structure: `<type>(<scope>): <subject>`
  - Types: feat, fix, docs, style, refactor, test, chore
  - Example: `feat(auth): add OAuth2 authentication`

## Merge Workflow
- Always create pull requests for merging into main/master
- Require at least one code review before merging
- Squash commits when merging feature branches
- Delete branches after merging
- Keep pull requests focused on a single feature/fix

## Git Actions
When asked to help with git operations:
- Create new branches from the latest main/master
- Stage only relevant files for commits
- Verify changes before committing
- Push to remote after commits
- Create pull requests with descriptive titles and descriptions
- Resolve merge conflicts by analyzing both versions carefully
- Rebase feature branches on main/master when needed
- Use git stash for temporary work storage

## Best Practices
- Commit early and often
- Write meaningful commit messages
- Never commit sensitive information
- Keep branches up to date with main/master
- Use .gitignore for excluding build artifacts and dependencies