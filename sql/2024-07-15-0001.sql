CREATE SCHEMA groat_pay;

CREATE TYPE groat_pay.KindOfPayment AS ENUM ('INCOMING', 'OUTGOING', 'INTERNAL');
CREATE TYPE groat_pay.KindOfPaymentState AS ENUM ('CREATED', 'DECLINED', 'SUCCEEDED', 'PROCESSING', 'DELAYED', 'ABORTED');
CREATE TYPE groat_pay.KindOfCurrency AS ENUM ('USD', 'EUR');

CREATE TABLE groat_pay.payments(
                                  id uuid NOT NULL ,
                                  user_id text NOT NULL,
                                  amount int NOT NULL ,
                                  kind groat_pay.KindOfPayment NOT NULL,
                                  currency groat_pay.KindOfCurrency NOT NULL DEFAULT 'USD',
                                  state groat_pay.KindOfPaymentState NOT NULL DEFAULT 'CREATED',
                                  created_at timestamp with time zone NOT NULL DEFAULT now(),
                                  processed_at timestamp with time zone
);
