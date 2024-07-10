import os
import json
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from google.oauth2 import service_account
from google.cloud import dialogflow

@csrf_exempt
def dialogflow_session(request):
    if request.method == 'POST':
        # Load JSON data from request body
        data = json.loads(request.body)
        session_id = data.get('session_id', 'default_session')
        text = data.get('text', 'Hello')

        # Set up Dialogflow session
        project_id = os.getenv('DIALOGFLOW_PROJECT_ID')
        session = f"projects/{project_id}/agent/sessions/{session_id}"
        credentials_json = json.loads(os.getenv('GOOGLE_CREDENTIALS_JSON'))
        credentials = service_account.Credentials.from_service_account_info(credentials_json)

        # Create a Dialogflow client
        session_client = dialogflow.SessionsClient(credentials=credentials)
        text_input = dialogflow.TextInput(text=text, language_code="en-US")
        query_input = dialogflow.QueryInput(text=text_input)
        response = session_client.detect_intent(request={'session': session, 'query_input': query_input})

        # Send response
        return JsonResponse({
            "fulfillmentText": response.query_result.fulfillment_text
        })

    return JsonResponse({"error": "This endpoint supports only POST method."})