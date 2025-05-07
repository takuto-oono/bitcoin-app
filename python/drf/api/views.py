from rest_framework import viewsets
from .models.ticker import Ticker
from .serializers.ticker import TickerSerializer


class TickerViewSet(viewsets.ModelViewSet):
    queryset = Ticker.objects.all()
    serializer_class = TickerSerializer
