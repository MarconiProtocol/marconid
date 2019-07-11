# Contribution Guidelines

## Overview
As a completely open source project, we welcome your contributions to Marconi Protocol and have put together a few guidelines to help you participate

## How to contribute (in a nutshell)
- You will need a Github account
- Fork the repository
- Create topic branch
- Commit patches

## Getting started
- Find an issue to tackle from our list of open issues here. Please only work on issues that have been assigned a priority. We also have tagged issues that are good for [first time contributors.]
- If you don’t see an issue you’d like to work on,  create a Github issue providing a detailed description and repro steps if you’re fixing a bug.  Once we have reviewed the issue internally and assigned a priority, you are free to work on the issue.
- We are focusing on bugs at this time; if you have feature requests, please tag them and we’ll add it to our feature pipeline for triaging.

## Development Process
### Making changes
- Create a topic branch from where you want to base your work.
  - This is usually the master branch.
  - Please avoid working directly on the master branch.
- Make sure you have added the necessary tests for your changes and make sure all tests pass.
### Submitting changes
- Push changes to a topic branch in your fork of the repository.
- Submit a pull request to the repository for the project you’re working on Include a descriptive commit message.
- Changes contributed via pull request should focus on a single issue at a time.
- Rebase your local changes against the master branch. Resolve any conflicts that arise.
- We will comment on pull requests within 1-2 business days and may provide feedback and suggestions.

### Coding style
GoLang projects: marconid, marconi_cli, go-ethereum
https://golang.org/doc/effective_go.html

Solidity projects: smart-contract-network
https://solidity.readthedocs.io/en/v0.5.3/style-guide.html

Javascript projects: middleware
https://google.github.io/styleguide/jsguide.html
