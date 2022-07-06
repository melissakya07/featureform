from featureform import ResourceClient
import featureform as ff

# rc = ResourceClient("localhost:8000")
redis = ff.get_redis("redis-quickstart")
postgres = ff.get_postgres("postgres-quickstart")

ff.register_user("featureformer").make_default_owner()

transactions = postgres.register_table(
    name="transactions",
    variant="kaggle",
    description="Fraud Dataset From Kaggle",
    table="Transactions",  # This is the table's name in Postgres
)


@postgres.sql_transformation(variant="quickstart")
def average_user_transaction_v2():
    """the average transaction amount for a user """
    return "SELECT CustomerID as user_id, avg(TransactionAmount) " \
           "as avg_transaction_amt from {{transactions.kaggle}} GROUP BY user_id"


user = ff.register_entity("user")

average_user_transaction_v2.register_resources(
    entity=user,
    entity_column="user_id",
    inference_store=redis,
    features=[
        {"name": "avg_transactions", "variant": "quickstart", "column": "avg_transaction_amt", "type": "float32"}
    ]
)

# Register label from our base Transactions table
transactions.register_resources(
    entity=user,
    entity_column="customerid",
    labels=[
        {"name": "fraudulent", "variant": "quickstart", "column": "isfraud", "type": "bool"}
    ],
)

ff.register_training_set(
    "fraud_training_v2", "quickstart",
    label=("fraudulent", "quickstart"),
    features=[("avg_transactions", "quickstart")],
)
