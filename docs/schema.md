# Database Schema Documentation

## Overview
This document describes the database schema for the Book Reader application. The schema is designed to manage book sources, book content, chapters, user favorites, and categorization.

## Tables

### source
Stores information about book source websites.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Unique identifier for the source |
| name | TEXT | NOT NULL, UNIQUE | Name of the book source |
| domain | TEXT | NOT NULL, UNIQUE | Domain URL of the source website |
| icon_url | TEXT | | URL to the source's icon/logo |
| updated_at | DATETIME | NOT NULL | Timestamp of last update |
| deleted_at | DATETIME | | Timestamp of soft deletion (NULL if not deleted) |

### book
Stores book metadata and information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Unique identifier for the book |
| url | TEXT | NOT NULL | URL to the book page |
| title | TEXT | NOT NULL | Title of the book |
| type | TEXT | NOT NULL, CHECK(type IN ('manga', 'novel')) | Type of book content (manga or novel) |
| image_url | TEXT | | URL to the book cover image |
| synopsis | TEXT | | Description/summary of the book |
| slug | TEXT | | URL-friendly identifier |
| other_id | TEXT | | External identifier from source |
| source_id | INTEGER | | Foreign key reference to source.id |
| updated_at | DATETIME | NOT NULL | Timestamp of last update |
| deleted_at | DATETIME | | Timestamp of soft deletion (NULL if not deleted) |

**Indexes:**
- `book_title_unique`: Unique index on (source_id, title) - ensures book titles are unique per source

### chapter
Stores individual book chapters.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Unique identifier for the chapter |
| url | TEXT | NOT NULL | URL to the chapter page |
| title | TEXT | NOT NULL | Title of the chapter |
| book_id | INTEGER | | Foreign key reference to book.id |
| number | TEXT | | Chapter number (stored as text to support formats like "1.5") |
| content | TEXT | | Serialized chapter content (image URLs for manga, text for novels) |
| upload_date | DATETIME | | Date the chapter was uploaded to the source |
| completed | BOOLEAN | NOT NULL, DEFAULT FALSE | Whether the chapter reading is completed |
| downloaded | BOOLEAN | NOT NULL, DEFAULT FALSE | Whether the chapter images are downloaded |
| other_id | TEXT | | External identifier from source |
| updated_at | DATETIME | NOT NULL | Timestamp of last update |
| deleted_at | DATETIME | | Timestamp of soft deletion (NULL if not deleted) |

**Indexes:**
- `chapter_number_unique`: Unique index on (book_id, number) - ensures chapter numbers are unique per book

### user
Stores user information.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Unique identifier for the user |
| name | TEXT | | User's name |
| updated_at | DATETIME | NOT NULL | Timestamp of last update |
| deleted_at | DATETIME | | Timestamp of soft deletion (NULL if not deleted) |

### favorite
Tracks user favorites and reading progress.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Unique identifier for the favorite |
| user_id | INTEGER | | Foreign key reference to user.id |
| book_id | INTEGER | | Foreign key reference to book.id |
| progress | TEXT | | Reading progress information |
| categories | TEXT | | Serialized list of category assignments |
| updated_at | DATETIME | NOT NULL | Timestamp of last update |
| deleted_at | DATETIME | | Timestamp of soft deletion (NULL if not deleted) |

### category
Stores user-defined categories for organizing books.

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| id | INTEGER | PRIMARY KEY, AUTOINCREMENT | Unique identifier for the category |
| name | TEXT | | Category name |
| updated_at | DATETIME | NOT NULL | Timestamp of last update |
| deleted_at | DATETIME | | Timestamp of soft deletion (NULL if not deleted) |

## Relationships

```
source (1) ----< (N) book
book (1) ----< (N) chapter
user (1) ----< (N) favorite
book (1) ----< (N) favorite
```

## Design Notes

1. **Timestamps**: The `updated_at`, `deleted_at`, and `upload_date` fields use DATETIME type to store timestamps in ISO 8601 format.

2. **Serialized Data**: Some fields like `content` and `categories` store serialized data as TEXT, suggesting JSON or comma-separated format.

3. **Chapter Numbers**: Chapter numbers are stored as TEXT to support decimal chapter numbers (e.g., "1.5", "2.1") and special formats.

4. **Unique Constraints**: 
   - Book titles are unique per source
   - Chapter numbers are unique per book
   - Source names and domains are globally unique

5. **Boolean Flags**: The chapter table uses boolean flags (`completed`, `downloaded`) to track reading and download status.

6. **External IDs**: The `other_id` fields allow storing identifiers from external source systems for synchronization purposes.

7. **Soft Deletes**: All tables include a `deleted_at` field for soft deletion support. Records with a non-NULL `deleted_at` value are considered deleted but remain in the database for audit purposes.
