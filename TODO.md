# dailytracker tasks

Never update this file or any of the files in the `/tasks`.

## 📋 To Do

### Sprint 2: 01.05.2026 - 22.05.2026
- [ ] [DT-1][Technical] Migrate to Postgres (local = sqlite, prod = Postres)
- [ ] [DT-2][Technical] Add roles (admin/user)
- [ ] [DT-3][Technical] Add some SQL library to not write raw SQL everywhere

### Backlog
- [ ] [Admin] Add admin user panel
- [ ] [Admin] Admin can add new users
- [ ] [Technical] Add better architecture - events and packages entries and statistics
- [ ] [Statistics] Add weekly/monthly statistics view
- [ ] [Statistics] Add data visualization (charts)
- [ ] [Admin/User] Add "Request user account" functionality
- [ ] [UI] Add dark mode toggle

### 🐛 Bugs

- [ ] (none currently)

### 💡 Ideas / Future Features

- [ ] Integration with calendar apps
- [ ] Streak tracking
- [ ] Goal setting and tracking
- [ ] Customizable scoring ranges
- [ ] Import data from other trackers
- [ ] nginx on server
- [ ] Add backup automation script
- [ ] Add data export (CSV/JSON)
- [ ] Add entry notes/comments field
- [ ] Add tags/categories for entries
- [ ] Add badges :D

### ✅ Done

- [x] Switch to pure Go SQLite (modernc.org/sqlite)
- [x] Create Docker setup
- [x] Add VPS deployment script
- [x] Update documentation
- [x] Fix database initialization bug
- [x] Add comprehensive .gitignore

#### Sprint 1: 14.04.2026 - 30.04.2026, PR: #1
- [x] [Technical] Update README.md
- [x] [Users] Add user authentication
- [x] [Technical] Change project structure to fit [this standard](https://github.com/golang-standards/project-layout?tab=readme-ov-file)
- [x] [Technical] Add some migration logic based on files
- [x] [Technical] Eager DB (or at least migration)
