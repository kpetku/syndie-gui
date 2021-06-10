# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
### Added
- Add go mod
- Add avatars to the message thread list
- Add CHANGELOG.md
- Add application versioning
- Rehash GUI after importing messages
- Add latest message view
- Add anonymous/non-anonymous fetch
- Inline image support
- Add Edit -> Settings menu

### Changed
- Updated GUI toolkit to Fyne v2.0.3
- Fix issue with messages containing long subject lines not wrapping
- Darkmode style fixes
- Fix stretched avatars in channel list
- Fixed file permssions on Linux
- Reworked file/folder/remote URL fetching dialogs
- Fixed "jumpy" scroll bars while changing channels or messages
- Fixed issue with some messages not being tappable
- Change the message view to a more modern card layout
- Cache avatars in memory
- Major UI rework

## [0.0.2] - 2020-07-31
### Added
- Add fetch from directory

### Changed
- Make fetch URL configurable

## [0.0.1] - 2020-07-25
### Added
- Initial preview release
