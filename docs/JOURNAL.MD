# Journal

## March 20, 2023



## December 18, 2022

I've researched a bit on GitHub or git-based integrations for read/write. I think that it would be best if I avoid that market because I don't want to be focused on that kind of business (UI, interactions, etc).

If ever, I could possibly see myself building an integration tool on top of platforms that are focused on that, but not build the whole infrastructure for the git integration itself.

I did a bit of research and found that [Gitbook](https://www.gitbook.com/) which I already use for tome.gg itself is a good git-related integration tool.

So if I work on the assumption that my business is not solely focused on implementing git-based integration simplified for non-developers (which GitBook itself deals with), then what does tome.gg focus on?

Its focus is the definition of its own protocol. GitBook can be the integration medium, but tome.gg can be built such that it is functional independent of the integration medium.

### Offloading work

Capabilities of the librarian protocol delegated to the integration medium:

- Data loading (Gitbook already has GitSync)
- Repository access grants (Gitbook has a GitHub or GitLab integration)
- User account identification (Gitbook has their own authentication, and also integrates the GitHub or GitLab OAuth access to the repositories)

Capabilities of the librarian protocol that we can and should focus on:

- Data validation (validating that a Git repository successfully matches the rules for the Librarian protocol)
- Proxied operations (initial set up of a Git repository)

### Workflow

Possibly, what we can do is initiate building the GitBook website, and then building a CLI tool that works on top of it.

- Meta information can be kept at a hidden directory called `.tome/`.
- Tome related configuration can be found in `tome.yaml`.