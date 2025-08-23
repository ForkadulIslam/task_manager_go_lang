# Backend Summary: Task Manager (Golang + Gin)

This document summarizes the current state and key features of the Task Manager backend, developed using Golang and the Gin framework.

## 1. Tech Stack
- **Backend Framework:** Gin (Golang)
- **Database:** MySQL
- **ORM:** GORM
- **Authentication:** JWT (JSON Web Tokens)

## 2. Core Features Implemented
- **User Management:** Registration, Login.
- **Group Management:** Full CRUD operations for groups.
- **Task Type Management:** Full CRUD operations for task types.
- **Task Management:**
    - Creation, retrieval (all, by ID), update, and deletion of tasks.
    - Assignment of tasks to individual users and groups.
    - Assignment of follow-up users to tasks.
    - **Task Status Updates:** Assigned users (direct or via group) can update task status.
    - **Task Commenting:** Assigned users (direct or via group) and follow-up users can add comments.

## 3. Key API Endpoints

### Authentication
- `POST /register`: Register a new user.
- `POST /login`: Authenticate a user and receive a JWT token.

### Tasks
- `POST /tasks`: Create a new task.
- `GET /tasks`: Retrieve all tasks.
- `GET /tasks/:id`: Retrieve a single task by ID. **Includes preloaded assigned users, assigned groups, follow-up users, and comment history (with user details).**
- `PUT /tasks/:id`: Update an existing task (authorized for task creator).
- `DELETE /tasks/:id`: Delete a task (authorized for task creator).
- `POST /tasks/:id/status`: Update the status of a task. (Authorized for assigned users - direct or via group).
- `POST /tasks/:id/comments`: Add a comment to a task. (Authorized for assigned users - direct or via group, and follow-up users).

### Groups
- `POST /groups`: Create a new group.
- `GET /groups`: Retrieve all groups.
- `GET /groups/:id`: Retrieve a single group by ID.
- `PUT /groups/:id`: Update an existing group.
- `DELETE /groups/:id`: Delete a group.

### Task Types
- `POST /task-types`: Create a new task type.
- `GET /task-types`: Retrieve all task types.
- `GET /task-types/:id`: Retrieve a single task type by ID.
- `PUT /task-types/:id`: Update an existing task type.
- `DELETE /task-types/:id`: Delete a task type.

## 4. Authorization Logic
- **Task Creator:** Has full CRUD access to tasks they created.
- **Assigned User (Direct or via Group):** Can update task status and add comments to assigned tasks.
- **Follow-up User:** Can only add comments to tasks they are following.
- All authenticated routes are protected by JWT middleware.

## 5. Database Schema Highlights
- `users`: Stores user credentials and roles.
- `groups`: Stores group information.
- `user_groups`: Junction table for many-to-many relationship between users and groups.
- `task_types`: Stores predefined task categories.
- `tasks`: Main task details, linked to `task_types` and `users` (for `CreatedBy`).
- `assign_task_to_users`: Links tasks to directly assigned users.
- `assign_task_to_groups`: Links tasks to assigned groups.
- `task_followup_users`: Links tasks to users designated for follow-up.
- `task_status_update_logs`: Logs all status changes for tasks, including who changed it and when.
- `task_comment_logs`: Stores comments made on tasks, linked to the user who made the comment.

## 6. Current State

The backend is fully functional and ready for integration with the frontend application. All planned API endpoints are implemented, and authorization rules are enforced.
