# Code Documentation Rule

## General Documentation Guidelines
- Document all public APIs, classes, and functions
- Keep comments up-to-date with code changes
- Use clear, concise language
- Focus on why and what, not how (unless complex)
- Document assumptions and edge cases

## Code Comments
- Add header comments for all files explaining purpose and usage
- Comment complex algorithms and business logic
- Explain "why" rather than "what" the code does
- Use TODO/FIXME comments for temporary solutions with issue references
- Keep comments professional and focused on technical aspects
- Remove commented-out code; use version control instead

## Comment Style by Language
- **Python**: Use docstrings (""") for modules, classes, methods, and functions
- **JavaScript/TypeScript**: Use JSDoc format for functions and classes
- **Java/C#**: Use standard documentation comments (/** */)
- **SQL**: Comment complex queries explaining the purpose and expected results
- **YAML/JSON**: Add comments for configuration files explaining purpose of sections

## README Files
- Every project must have a README.md in the root directory containing:
  - Project name and brief description
  - Purpose and problem it solves
  - Installation instructions
  - Usage examples
  - Key dependencies and versions
  - Configuration requirements
  - Contribution guidelines
  - License information

## Module Documentation
- Each major module/component should have its own README.md explaining:
  - Module's purpose and functionality
  - How it integrates with other components
  - API documentation or usage examples
  - Any specific configuration required

## API Documentation
- Document all endpoints with:
  - HTTP method and URL
  - Request parameters and body format
  - Response format and status codes
  - Authentication requirements
  - Rate limiting information
  - Example requests and responses

## Documentation Actions
When asked to help with code:
- Add appropriate file header comments
- Include docstrings/comments for all functions and classes
- Create or update README.md files for new components
- Document dependencies and their versions
- Explain complex logic with inline comments
- Ensure comments follow the project's existing style