from rest_framework import viewsets
from rest_framework.decorators import api_view
from rest_framework.response import Response
from .models.ticker import Ticker
from .serializers.ticker import TickerSerializer


@api_view(['GET'])
def healthcheck(request):
    return Response("ok")


class TickerViewSet(viewsets.ModelViewSet):
    queryset = Ticker.objects.all()
    serializer_class = TickerSerializer
