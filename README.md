# ai-chatbot-backend

## python-chatbot

## Dialogflow Integration Service

This project is a Django-based web service that integrates with Google Dialogflow to process and respond to user messages within a session context.

## Prerequisites

Before you begin, ensure you have met the following requirements:
* You have installed Python 3 or higher.
* You have a Django-supported database installed (SQLite is used in this example).
* You have an active Google Cloud account with access to Dialogflow.

```bash
# cd into your project path
python -m venv myvenv
source myvenv/bin/activate  # On Windows use `myvenv\Scripts\activate`
pip install -r requirements.txt 
```
Create a .env file in the project root directory and include necessary configurations:
```
GOOGLE_CREDENTIALS_JSON=<your-credentials-json>
DIALOGFLOW_PROJECT_ID=<your-project-id>
```
Apply migrations to set up the database schema:
```
python manage.py migrate
```
Launch the development server:
```
python manage.py runserver
```