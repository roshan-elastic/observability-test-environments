# End Of Life (EOL)

We are committed to providing support for our tools.
We make frequent releases to include new features, bug fixes, and security updates.

Oblt-cli use semantic versioning, Major.Minor.Patch.
Minor and Patch releases are backward compatible, so no breaking changes are introduced.
Major releases may introduce breaking changes, and we will provide migration guides and deprecation warnings.
Every time we release a major version, we stop supporting the previous major version, and the update is mandatory.
Users that try to use oblt-cli in a previous major version will receive an error message and a link to the migration guide, the tool stop working for them until they update.

The reason for this is to keep the codebase clean and maintainable, and to avoid technical debt.
Our tooling is focused on our daily work, and the update to the new version is simple and fast.
There is no reason that blocks developers to update to the latest version.
