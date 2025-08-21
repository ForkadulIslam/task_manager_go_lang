📌 Task Manager Backend (Golang + Gin)
Tech Stack

Backend: Golang (Gin framework)

Frontend: Vue.js

Database: MySQL

ORM: GORM (Laravel-Eloquent-like ORM in Golang)

Project Story

A Task Manager application with users, groups, tasks, and follow-up functionality.

Authentication & Access

Users log in with username + password.

JWT-based authentication for session handling.

User roles:

Super Admin → Full access (can manage users, groups, and all tasks).

User → Access limited to tasks created by them or assigned to them.

Task Features

A user can create a new task with:

Task type, priority, start date, due date, description, attachments.

Task assignment:

Assign to individual users (from user dropdown).

Assign to groups (from group dropdown → applies to all members).

Follow-up assignment:

A task may also be assigned to follow-up users/groups.

Follow-up users can only add remarks to the task.

Task access rules:

Task creator → full CRUD.

Assigned user → can update status & progress notes.

Follow-up user → can only add remarks.

Task statuses: Pending, In Progress, In Review, Completed.

Status changes & comments are logged.

Group Management

A group has:

Name (label).

One or many users.

One user can belong to multiple groups.

Groups are created by users (creator is tracked).

Database Schema
users

id (PK)

username (varchar, unique)

password (hashed)

status (tinyint → 1=Active, 0=Inactive)

user_label (enum → 1=Super Admin, 2=User)

created_at, updated_at

groups

id (PK)

label (varchar)

created_by (FK → users.id)

created_at, updated_at

user_groups

id (PK)

user_id (FK → users.id)

group_id (FK → groups.id)

task_types

id (PK)

label (varchar)

created_at, updated_at

tasks

id (PK)

label (varchar)

task_type_id (FK → task_types.id)

priority (enum: Normal, Medium, High, Escalation)

start_date (datetime)

due_date (datetime)

description (longtext)

attachment (varchar, nullable)

status (enum: 0=Pending, 1=In Progress, 2=In Review, 3=Completed)

created_by (FK → users.id)

created_at, updated_at

assign_task_to_users

id (PK)

user_id (FK → users.id)

task_id (FK → tasks.id)

assign_task_to_groups

id (PK)

group_id (FK → groups.id)

task_id (FK → tasks.id)

task_followup_users

id (PK)

user_id (FK → users.id)

task_id (FK → tasks.id)

remarks (text, nullable)

created_at, updated_at

task_status_update_logs

id (PK)

task_id (FK → tasks.id)

user_id (FK → users.id)

status (enum: 0=Pending, 1=In Progress, 2=In Review, 3=Completed)

created_at, updated_at

task_comment_logs

id (PK)

task_id (FK → tasks.id)

user_id (FK → users.id)

comment (text)

created_at, updated_at

task_seen_by_users

id (PK)

task_id (FK → tasks.id)

user_id (FK → users.id)

created_at, updated_at