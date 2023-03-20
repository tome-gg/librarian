# System Capabilities

1. **Data validation** - Given a Github repository, the system should be able to run validation checks: repository + protocol version = errors/warnings.

This enables us both the mentor and the apprentice to double check their entry, to ensure that the data can be loaded and processed by Tome.gg's librarian.

Part of validation should be to (a) ensure that all trainings have been evaluated.

# Deferred system capabilities

## Nearby

1. **User account identification** - a website visitor should be able to authenticate using a Github account to identify themselves. 

2. **Repository access grants** - an authenticated user should be able to grant read and write access to their own set of repositories. This will enable tome.gg to write journal entries in behalf of the user.

3. **Data loading** - the system should be able to load a an accessible Github repository (whether public, or authenticated access) to load data into tome.gg's smart contracts.

## Far into the future

1. **Proxied operations** - Creating a new repository for an authenticated user (e.g. for their journaling, or for DSU purposes), or granting or revoking access to repositories. If we can do this on tome.gg's application itself, that would be more convenient and less fragmented.