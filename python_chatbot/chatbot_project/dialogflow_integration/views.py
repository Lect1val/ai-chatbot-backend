import os
import json
import openai
from django.http import JsonResponse
from django.views.decorators.csrf import csrf_exempt
from google.oauth2 import service_account
from google.cloud import dialogflow

@csrf_exempt
def dialogflow_session(request):
    if request.method == 'POST':
        data = json.loads(request.body)
        project_id = os.getenv('DIALOGFLOW_PROJECT_ID')
        session_id = data.get('session_id', 'default_session')
        text = data.get('text', 'Hello')

        session_path = f"projects/{project_id}/agent/sessions/{session_id}"
        credentials_json = json.loads(os.getenv('GOOGLE_CREDENTIALS_JSON'))
        credentials = service_account.Credentials.from_service_account_info(credentials_json)

        session_client = dialogflow.SessionsClient(credentials=credentials)
        text_input = dialogflow.TextInput(text=text, language_code="en-US")
        query_input = dialogflow.QueryInput(text=text_input)
        response = session_client.detect_intent(request={'session': session_path, 'query_input': query_input})

        print("Full Dialogflow Response:", response)

        intent_name = response.query_result.intent.display_name

        print("Intent Name:", intent_name)

        # Directly access parameters without merging
        session_name = response.query_result.parameters.get('session_name', None)
        print("Session Name:", session_name)

        if intent_name == "session search":
            if session_name:
                return JsonResponse({"fulfillmentText": f"Information for session: {session_name}"})
            else:
                return JsonResponse({"fulfillmentText": "Session name not provided."})

        elif intent_name == "Aster arcade URL":
            return JsonResponse({
                "fulfillmentText": "This is the URL of Aster Arcade: [Aster-arcade](url:https://aster.arisetech.dev/aster-arcade/)"
            })

        if response.query_result.intent.is_fallback or response.query_result.intent.display_name == "":
            # Fallback to OpenAI using the new API
            openai.api_key = os.getenv('OPENAI_API_KEY')
            openai_response = openai.ChatCompletion.create(
                model="gpt-3.5-turbo",
                messages=[{"role": "user", "content": text}]
            )
            return JsonResponse({"fulfillmentText": openai_response['choices'][0]['message']['content'].strip()})
        else:
            return JsonResponse({"fulfillmentText": response.query_result.fulfillment_text})

        # Default response if no specific intent matched
        # return JsonResponse({"fulfillmentText": response.query_result.fulfillment_text})

    return JsonResponse({"status": "error", "message": "This endpoint supports only POST method."})
