# Reflection on Building with AI

Building the Go Image Processing Service was a fascinating exercise in human-AI collaboration. The process moved at a remarkable pace, largely due to a structured, iterative workflow that leveraged the AI's strengths in code generation, testing, and rapid prototyping. This reflection explores how that partnership impacted the build, what worked well, what felt limiting, and the key lessons learned about prompting, reviewing, and iterating with an AI partner.

## The Impact of AI on the Build Process

The most significant impact of using an AI assistant was the dramatic acceleration of the development cycle. What would have been a multi-hour or multi-day project for a solo developer learning the ropes was condensed into a single, highly productive session. The core of this efficiency came from offloading the cognitive burden of boilerplate and syntax. Instead of getting bogged down in the minutiae of setting up an HTTP server, parsing multipart forms, or remembering the exact function signature for a library, I could focus on the higher-level architectural decisions: what feature to build next, what the API should look like, and what edge cases to consider.

The AI acted as a tireless pair programmer. It generated not just the implementation but also the corresponding test cases, enforcing a Test-Driven Development (TDD) methodology from the outset. This was a critical success factor. The TDD loop—prompting for tests first (Red), then prompting for the implementation to make them pass (Green), and finally refactoring—created a robust and predictable rhythm. It built a safety net that allowed for fearless refactoring and ensured that each new feature was solid before moving on.

## What Worked Well

1.  **Test-Driven Development (TDD) as a Guiding Framework**: The TDD approach was the perfect structure for this collaboration. It provided a clear, unambiguous definition of "done" for each feature. Prompting "write the tests for a `/crop` endpoint" created a concrete contract that the AI could then fulfill with its implementation. This eliminated ambiguity and made the process highly efficient.

2.  **Iterative Prompting and Debugging**: The AI wasn't perfect, and this was a good thing. Early on, it produced code that had subtle bugs, such as forgetting to rewind a file handle after parsing a form or omitting a necessary package import for PNG decoding. This created a natural, collaborative debugging loop. I would run the code, paste the exact error message back to the AI, and it would immediately understand the context, explain the error, and provide the corrected code. This call-and-response felt less like dealing with a faulty tool and more like working with a junior developer who occasionally needs guidance.

3.  **High-Level Brainstorming and Planning**: At the start and before each new feature, I could ask the AI for a high-level approach. For example, when considering PDF compression, the AI correctly identified that the existing libraries were unsuitable and proposed a new strategy involving a command-line tool (`ghostscript`). This ability to reason about the problem space and suggest appropriate tools was invaluable.

## What Felt Limiting

The primary limitation was the AI's lack of persistent state and its reliance on the immediate context provided. On a few occasions, the `replace` tool failed because the `old_string` I provided wasn't a perfect, character-for-character match of the file's current content. This required an extra step of reading the file first to get the exact text before attempting a modification. A human developer would intuitively remember the changes they just made, but the AI required this explicit context to be re-established for each atomic operation. While the tooling provides a good abstraction, this underlying statelessness is a friction point that requires conscious management.

## Lessons Learned About Prompting, Reviewing, and Iterating

This project crystallized several key lessons about working effectively with AI:

-   **Be the Architect, Let the AI Be the Builder**: My most effective role was to set the direction. By defining the API contract (`/crop` needs `x, y, width, height`), the acceptance criteria (the TDD tests), and the overall structure, I could steer the project effectively. The AI excelled at the tactical work of filling in the implementation details.

-   **Errors Are Your Best Prompts**: The most effective prompts were often the direct error messages from the Go compiler or the test runner. Instead of trying to describe the problem in natural language, providing the raw, technical output gave the AI the precise information it needed to identify and fix the issue.

-   **Review Everything, Trust but Verify**: The AI-generated code was generally high-quality, but it was not infallible. Reviewing every line of code was non-negotiable. This was essential for catching the subtle bugs mentioned earlier and for ensuring the code adhered to the project's evolving conventions. The goal is not to blindly accept the AI's output but to use it as a high-quality first draft that you, the developer, are responsible for validating.

In conclusion, the partnership was a resounding success. By embracing a structured TDD workflow and understanding the AI's strengths and limitations, we were able to build a functional and well-tested application at a speed that would be impossible to achieve alone. The process felt less like "using a tool" and more like a genuine collaboration, where my role was to provide the vision, strategy, and critical review, while the AI provided the tireless, high-speed implementation.
