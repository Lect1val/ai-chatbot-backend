from django.urls import path
from . import views

urlpatterns = [
    path('session/', views.dialogflow_session, name='dialogflow_session'),
]

