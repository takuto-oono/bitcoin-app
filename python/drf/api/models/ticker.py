from django.db import models


class Ticker(models.Model):
    """
    BitFlyerからのティッカー情報を格納するモデル
    """
    id = models.AutoField(primary_key=True)
    tick_id = models.IntegerField(unique=True)
    product_code = models.CharField(max_length=50)
    state = models.CharField(max_length=50)
    timestamp = models.CharField(max_length=50)
    best_bid = models.FloatField()
    best_ask = models.FloatField()
    best_bid_size = models.FloatField()
    best_ask_size = models.FloatField()
    total_bid_depth = models.FloatField()
    total_ask_depth = models.FloatField()
    market_bid_size = models.FloatField()
    market_ask_size = models.FloatField()
    ltp = models.FloatField()
    volume = models.FloatField()
    volume_by_product = models.FloatField()

    class Meta:
        db_table = 'tickers'
