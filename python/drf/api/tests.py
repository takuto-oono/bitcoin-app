from django.test import TestCase
from django.urls import reverse
from rest_framework import status
from rest_framework.test import APITestCase, APIClient
from .models.ticker import Ticker
from .serializers.ticker import TickerSerializer


class TickerSerializerTestCase(TestCase):
    """TickerSerializerのテストクラス"""

    def setUp(self):
        """テストデータのセットアップ"""
        self.ticker_data = {
            'tick_id': 12345,
            'product_code': 'BTC_JPY',
            'state': 'RUNNING',
            'timestamp': '2025-05-07T10:54:00',
            'best_bid': 5000000.0,
            'best_ask': 5000100.0,
            'best_bid_size': 0.1,
            'best_ask_size': 0.2,
            'total_bid_depth': 100.0,
            'total_ask_depth': 200.0,
            'market_bid_size': 300.0,
            'market_ask_size': 400.0,
            'ltp': 5000050.0,
            'volume': 500.0,
            'volume_by_product': 600.0
        }
        self.ticker = Ticker.objects.create(**self.ticker_data)
        self.serializer = TickerSerializer(instance=self.ticker)

    def test_serializer_contains_expected_fields(self):
        """シリアライザーが期待されるフィールドを含んでいるかテスト"""
        data = self.serializer.data
        expected_fields = {
            'id', 'tick_id', 'product_code', 'state', 'timestamp',
            'best_bid', 'best_ask', 'best_bid_size', 'best_ask_size',
            'total_bid_depth', 'total_ask_depth',
            'market_bid_size', 'market_ask_size',
            'ltp', 'volume',
            'volume_by_product'
        }
        self.assertEqual(set(data.keys()), expected_fields)

    def test_serializer_field_content(self):
        """シリアライザーのフィールド内容が正しいかテスト"""
        data = self.serializer.data
        self.assertEqual(data['tick_id'], self.ticker_data['tick_id'])
        self.assertEqual(data['product_code'],
                         self.ticker_data['product_code'])
        self.assertEqual(data['state'], self.ticker_data['state'])
        self.assertEqual(data['timestamp'], self.ticker_data['timestamp'])
        self.assertEqual(float(data['best_bid']), self.ticker_data['best_bid'])
        self.assertEqual(float(data['best_ask']), self.ticker_data['best_ask'])

    def test_deserialize_with_valid_data(self):
        """有効なデータでデシリアライズできるかテスト"""
        new_ticker_data = {
            'tick_id': 67890,
            'product_code': 'ETH_JPY',
            'state': 'RUNNING',
            'timestamp': '2025-05-07T11:00:00',
            'best_bid': 300000.0,
            'best_ask': 300100.0,
            'best_bid_size': 0.3,
            'best_ask_size': 0.4,
            'total_bid_depth': 150.0,
            'total_ask_depth': 250.0,
            'market_bid_size': 350.0,
            'market_ask_size': 450.0,
            'ltp': 300050.0,
            'volume': 550.0,
            'volume_by_product': 650.0
        }
        serializer = TickerSerializer(data=new_ticker_data)
        self.assertTrue(serializer.is_valid())
        ticker = serializer.save()
        self.assertEqual(ticker.tick_id, new_ticker_data['tick_id'])
        self.assertEqual(ticker.product_code, new_ticker_data['product_code'])

    def test_deserialize_with_invalid_data(self):
        """無効なデータでデシリアライズできないことをテスト"""
        # 既存のtick_idと重複（uniqueバリデーション違反）
        invalid_data = {
            'tick_id': 12345,
            'product_code': 'BTC_JPY',
            'state': 'RUNNING',
            'timestamp': '2025-05-07T10:54:00',
            'best_bid': 5000000.0,
            'best_ask': 5000100.0,
            'best_bid_size': 0.1,
            'best_ask_size': 0.2,
            'total_bid_depth': 100.0,
            'total_ask_depth': 200.0,
            'market_bid_size': 300.0,
            'market_ask_size': 400.0,
            'ltp': 5000050.0,
            'volume': 500.0,
            'volume_by_product': 600.0
        }
        serializer = TickerSerializer(data=invalid_data)
        self.assertFalse(serializer.is_valid())
        self.assertIn('tick_id', serializer.errors)


class TickerViewSetTestCase(APITestCase):
    """TickerViewSetのテストクラス"""

    def setUp(self):
        """テストデータのセットアップ"""
        self.client = APIClient()
        self.ticker_data = {
            'tick_id': 12345,
            'product_code': 'BTC_JPY',
            'state': 'RUNNING',
            'timestamp': '2025-05-07T10:54:00',
            'best_bid': 5000000.0,
            'best_ask': 5000100.0,
            'best_bid_size': 0.1,
            'best_ask_size': 0.2,
            'total_bid_depth': 100.0,
            'total_ask_depth': 200.0,
            'market_bid_size': 300.0,
            'market_ask_size': 400.0,
            'ltp': 5000050.0,
            'volume': 500.0,
            'volume_by_product': 600.0
        }
        self.ticker = Ticker.objects.create(**self.ticker_data)
        self.list_url = reverse('ticker-list')
        self.detail_url = reverse('ticker-detail', args=[self.ticker.id])

    def test_get_ticker_list(self):
        """ティッカーリストの取得テスト"""
        response = self.client.get(self.list_url)
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertEqual(len(response.data), 1)
        self.assertEqual(
            response.data[0]['tick_id'],
            self.ticker_data['tick_id']
        )

    def test_get_ticker_detail(self):
        """ティッカー詳細の取得テスト"""
        response = self.client.get(self.detail_url)
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.assertEqual(
            response.data['tick_id'],
            self.ticker_data['tick_id']
        )
        self.assertEqual(
            response.data['product_code'],
            self.ticker_data['product_code']
        )

    def test_create_ticker(self):
        """ティッカー作成テスト"""
        new_ticker_data = {
            'tick_id': 67890,
            'product_code': 'ETH_JPY',
            'state': 'RUNNING',
            'timestamp': '2025-05-07T11:00:00',
            'best_bid': 300000.0,
            'best_ask': 300100.0,
            'best_bid_size': 0.3,
            'best_ask_size': 0.4,
            'total_bid_depth': 150.0,
            'total_ask_depth': 250.0,
            'market_bid_size': 350.0,
            'market_ask_size': 450.0,
            'ltp': 300050.0,
            'volume': 550.0,
            'volume_by_product': 650.0
        }
        response = self.client.post(
            self.list_url,
            new_ticker_data,
            format='json'
        )
        self.assertEqual(response.status_code, status.HTTP_201_CREATED)
        self.assertEqual(Ticker.objects.count(), 2)
        new_ticker = Ticker.objects.get(tick_id=67890)
        self.assertEqual(new_ticker.product_code, 'ETH_JPY')

    def test_update_ticker(self):
        """ティッカー更新テスト"""
        updated_data = {
            'tick_id': 12345,
            'product_code': 'BTC_JPY',
            'state': 'STOPPED',  # 変更
            'timestamp': '2025-05-07T12:00:00',  # 変更
            'best_bid': 5100000.0,  # 変更
            'best_ask': 5100100.0,  # 変更
            'best_bid_size': 0.1,
            'best_ask_size': 0.2,
            'total_bid_depth': 100.0,
            'total_ask_depth': 200.0,
            'market_bid_size': 300.0,
            'market_ask_size': 400.0,
            'ltp': 5100050.0,  # 変更
            'volume': 500.0,
            'volume_by_product': 600.0
        }
        response = self.client.put(
            self.detail_url,
            updated_data,
            format='json'
        )
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.ticker.refresh_from_db()
        self.assertEqual(self.ticker.state, 'STOPPED')
        self.assertEqual(self.ticker.timestamp, '2025-05-07T12:00:00')
        self.assertEqual(self.ticker.best_bid, 5100000.0)

    def test_partial_update_ticker(self):
        """ティッカー部分更新テスト"""
        partial_data = {
            'state': 'PAUSED',
            'best_bid': 5200000.0,
        }
        response = self.client.patch(
            self.detail_url,
            partial_data,
            format='json'
        )
        self.assertEqual(response.status_code, status.HTTP_200_OK)
        self.ticker.refresh_from_db()
        self.assertEqual(self.ticker.state, 'PAUSED')
        self.assertEqual(self.ticker.best_bid, 5200000.0)
        # 他のフィールドは変更されていないことを確認
        self.assertEqual(self.ticker.product_code, 'BTC_JPY')
        self.assertEqual(self.ticker.tick_id, 12345)

    def test_delete_ticker(self):
        """ティッカー削除テスト"""
        response = self.client.delete(self.detail_url)
        self.assertEqual(response.status_code, status.HTTP_204_NO_CONTENT)
        self.assertEqual(Ticker.objects.count(), 0)
