from django.urls import include, path
from rest_framework import routers
from . import views

router = routers.SimpleRouter()
router.register('tickers', views.TickerViewSet)

urlpatterns = [
    path("healthcheck/", views.healthcheck, name="healthcheck"),
    path('', include(router.urls)),
]
