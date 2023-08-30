# Features:

## System

- Trending
- Category
- Tag
- News
- Best author
- Upload novel
- Novel:
    - Volume -> Chapter
    - Text
    - Image (Table)

- Search:
    - Magic search: Describe then find (LLM)

- Reading:
    - Bookmark
    - Usage Statistic
    - Subcribtion (Push notification)
    - Comment
    - Views
    - Rating:
        - Star based

- Account related stuff:
    - Name, Avatar, Up-post, Poll

- Text: Markdown

## Front-End

- Progress bar
- Font-size, font-family
- TTS
- Chapter
- 

## Auth
- JWT

## Database
User can report Contents
User can block other users
User can subcribe to other user
### Entity

- Users:
    - UUID
    - Avatar_url
    - Name
    - Email
    - NoNovels
    - Reported
    - ReportedFor

- Novels:
    - UUID
    - Title
    - Tagline
    - Description
    - FrontPage
    - User
    - Tags
    - Language
    - Adult
    - Status
    - CreatedAt
    - UpdatedAt
    - TotalRating
    - NumberOfRater
    - Views
    - Click
    - Images_ID
    - Reported
    - ReportedFor

- Tags:
    - UUID
    - Name
    - Description

- Volumes:
    - UUID
    - Title
    - Tagline
    - Description
    - FrontPage
    - CreatedAt
    - UpdatedAt
    - Views
    - Click
    - Reported
    - ReportedFor

- Chapters:
    - UUID
    - Title
    - CreatedAt
    - UpdatedAt
    - Contents
    - Views
    - Click
    - Reported
    - ReportedFor
    
- Images:
    - UUID
    - Name
    - URL 
    - UploadedAt

- Comments:
    - UUID
    - CommentTo
    - User
    - Content (Markdown)
    - Reported
    - ReportedFor

### Relation
| Entity   | Relation         | Entity   |Features    |
| Users    | Self-referencing | Users    |Block       |
| Users    | Self-referencing | Users    |Subcribtion |
| Users    | OneToMany        | Comments |Commenting  |
| Users    | OneToMany        | Novels   |Upload      |
| Novels   | OneToMany        | Images   |Upload      |
| Novels   | OneToMany        | Comments |Commenting  |
| Novels   | OneToMany        | Tags     |Tagging     |
| Novels   | OneToMany        | Volumes  |Upload      |
| Volumes  | OneToMany        | Comments | Commenting |
| Volumes  | OneToMany        | Chapters | Upload     |
| Chapters | OneToMany        | Comments | Commenting |
| Comments | Self-referencing | Comments | Commenting |
