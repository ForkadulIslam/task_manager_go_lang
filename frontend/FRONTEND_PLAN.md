# Frontend Development Plan for Task Manager

This document outlines the plan and architecture for the Vue.js frontend of the Task Manager application. It is based on the existing Golang backend API.

## 1. Tech Stack

-   **Framework**: Vue.js 3 (with Composition API)
-   **Routing**: Vue Router
-   **State Management**: Pinia
-   **HTTP Client**: Axios
-   **Styling**: Tailwind CSS (for rapid, utility-first styling)

## 2. Project Structure

A standard Vue 3 project structure will be used:

```
/frontend
├── public/
└── src/
    ├── assets/
    ├── components/      # Reusable UI components (e.g., TaskCard, AppButton)
    ├── views/           # Page components (e.g., Login, TaskList, TaskDetail)
    ├── services/        # API communication layer
    ├── stores/          # Pinia state management modules
    ├── router/          # Vue Router configuration
    └── main.js          # App entry point
```

## 3. API Service (`src/services/api.js`)

This service will manage all communication with the backend.

-   **Axios Instance**: A global Axios instance will be configured.
-   **Authorization Header**: An interceptor will be set up to automatically add the `Authorization: Bearer <token>` header to all outgoing requests, pulling the token from the Pinia auth store.
-   **API Functions**: The service will export functions for each backend endpoint:
    -   `login(credentials)`
    -   `register(userInfo)`
    -   `getTasks()`
    -   `getTaskById(id)`
    -   `createTask(taskData)`
    -   `updateTask(id, taskData)`
    -   `deleteTask(id)`
    -   `updateTaskStatus(id, status)`
    -   `addTaskComment(id, comment)`
    -   `getGroups()`, `getUsers()`, `getTaskTypes()`, etc.

## 4. State Management (`src/stores/`)

Pinia stores will manage the application's global state.

-   **`auth.js`**:
    -   **State**: `user`, `token`
    -   **Getters**: `isAuthenticated`
    -   **Actions**: `login()`, `logout()`, `register()`
-   **`tasks.js`**:
    -   **State**: `tasks` (array), `selectedTask` (object)
    -   **Actions**: `fetchTasks()`, `fetchTask(id)`, `createTask()`, `updateTask()`
-   **`meta.js`**:
    -   **State**: `users`, `groups`, `taskTypes`
    -   **Actions**: `fetchUsers()`, `fetchGroups()`, `fetchTaskTypes()` (to populate dropdowns in forms)

## 5. Views (Pages) (`src/views/`)

The application will consist of the following pages:

1.  **`Login.vue`**:
    -   A simple form for `username` and `password`.
    -   On successful login, it will save the token using the Pinia store and redirect to the dashboard/task list.
    -   Will display validation errors from the API.

2.  **`AppLayout.vue`**:
    -   A main layout component that includes the navigation bar and a `<router-view>` to render the pages. This will wrap all authenticated views.

3.  **`TaskList.vue`**:
    -   The default page after login.
    -   Fetches and displays a list of tasks using a `TaskCard` component for each.
    -   Includes a "Create New Task" button.
    -   Provides filtering and sorting options.

4.  **`TaskDetail.vue`**:
    -   Displays the full details of a single task, fetched by its ID.
    -   Shows `description`, `priority`, `status`, `dates`.
    -   Lists assigned users, groups, and follow-up users.
    -   Displays a log of comments (`TaskCommentLog`).
    -   Allows assigned users to update the task status.
    -   Allows assigned/follow-up users to add new comments.
    -   The task creator will see "Edit" and "Delete" buttons.

5.  **`TaskForm.vue`**:
    -   A single form for both creating and editing tasks.
    -   The form will be pre-filled with task data when editing.
    -   Will use dropdowns (populated from the `meta.js` store) for `Task Type`, `Assigned Users`, `Groups`, etc.

## 6. Routing (`src/router/index.js`)

-   **Routes**: Define paths for all views (`/login`, `/tasks`, `/tasks/:id`, `/tasks/new`, `/tasks/:id/edit`).
-   **Navigation Guard**: A global `beforeEach` guard will check `authStore.isAuthenticated`. If a route requires authentication and the user is not logged in, it will redirect to `/login`.

## 7. Backend Data Structures for Frontend

Based on the Go controllers, here are the expected JSON structures the frontend will work with.

### Task (as returned by `GET /tasks/:id`)

```json
{
  "ID": 1,
  "CreatedAt": "2023-10-27T10:00:00Z",
  "UpdatedAt": "2023-10-27T10:00:00Z",
  "DeletedAt": null,
  "Label": "Deploy to production",
  "TaskTypeID": 1,
  "Priority": "High",
  "StartDate": "2023-11-01T00:00:00Z",
  "DueDate": "2023-11-05T00:00:00Z",
  "Description": "Detailed description of the deployment task.",
  "Attachment": "/uploads/attachment.pdf",
  "Status": "In Progress",
  "CreatedBy": 2,
  "AssignedUsers": [
    {
      "ID": 1,
      "UserID": 3,
      "TaskID": 1,
      "User": { "ID": 3, "Username": "frontend_dev" }
    }
  ],
  "AssignedGroups": [
    {
      "ID": 1,
      "GroupID": 1,
      "TaskID": 1,
      "Group": {
        "ID": 1,
        "Label": "Backend Team",
        "Users": [
          { "ID": 4, "Username": "backend_dev1" },
          { "ID": 5, "Username": "backend_dev2" }
        ]
      }
    }
  ],
  "FollowupUsers": [
      {
          "ID": 1,
          "UserID": 6,
          "TaskID": 1,
          "User": { "ID": 6, "Username": "project_manager" }
      }
  ],
  "Comments": [
    {
      "ID": 1,
      "TaskID": 1,
      "UserID": 3,
      "Comment": "I've started working on the UI.",
      "CreatedAt": "2023-10-28T11:00:00Z",
      "User": { "ID": 3, "Username": "frontend_dev" }
    }
  ]
}
```
