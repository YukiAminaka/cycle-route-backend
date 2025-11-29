CREATE EXTENSION IF NOT EXISTS postgis;    -- PostGIS有効化
CREATE EXTENSION IF NOT EXISTS btree_gist; -- 便利（排他制約や複合インデックス用）
CREATE EXTENSION IF NOT EXISTS pg_trgm;    -- 名前/説明/検索用


CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,                    -- ユーザーID
    ulid  CHAR(26) NOT NULL UNIQUE,
    name TEXT NOT NULL,                          -- ユーザー名
    highlighted_photo_id BIGINT DEFAULT 0,       -- ハイライト写真ID
    locale VARCHAR(10) DEFAULT 'ja',             -- 言語設定
    created_at TIMESTAMPTZ NOT NULL,             -- 作成日時（タイムゾーン付き）
    description TEXT,                            -- 自己紹介
    locality TEXT,                               -- 地域（市区）
    administrative_area TEXT,                    -- 行政区（都道府県など）
    country_code CHAR(2) DEFAULT 'JP',           -- 国コード（ISO形式）
    postal_code         VARCHAR(20),             -- 郵便番号 (マップの中心座標)
    geom geometry(Point,4326),                   -- 位置
    first_name TEXT,                             -- 名
    last_name TEXT,                              -- 姓
    email TEXT UNIQUE,                           -- メールアドレス（ユニーク）
    has_set_location BOOLEAN DEFAULT FALSE       -- 位置情報設定済みフラグ
);

CREATE TABLE routes (
  id                  BIGSERIAL PRIMARY KEY,                 -- 例: 42125024
  ulid  CHAR(26) NOT NULL UNIQUE,
  user_id             BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name                TEXT NOT NULL,
  description         TEXT NOT NULL DEFAULT '',
  highlighted_photo_id        BIGINT      DEFAULT 0,
  has_course_points   BOOLEAN NOT NULL DEFAULT FALSE,
  distance            DOUBLE PRECISION NOT NULL CHECK (distance >= 0),   -- 距離(m)
  duration            INTEGER NOT NULL CHECK (duration IS NULL OR duration >= 0), -- 所要時間(s)
  elevation_gain      DOUBLE PRECISION NOT NULL DEFAULT 0 CHECK (elevation_gain >= 0),
  elevation_loss      DOUBLE PRECISION NOT NULL DEFAULT 0 CHECK (elevation_loss >= 0),
  path_geom           geometry(LineString, 4326) NOT NULL CHECK (NOT ST_IsEmpty(path_geom)) CHECK (ST_NPoints(path_geom) >= 2),  -- 経路パス 空ジオメトリや、点が1個だけの線を保存禁止
  bbox                geometry(Polygon,4326) NOT NULL GENERATED ALWAYS AS (ST_Envelope(path_geom)) STORED,    -- マップの表示領域 path_geomから自動生成
  first_point         geometry(Point,4326) NOT NULL GENERATED ALWAYS AS (ST_StartPoint(path_geom)) STORED,  -- スタート位置 path_geomの始点から自動生成
  last_point          geometry(Point,4326) NOT NULL GENERATED ALWAYS AS (ST_EndPoint(path_geom))   STORED, -- 終了位置　path_geomの終点から自動生成
  created_at          TIMESTAMPTZ NOT NULL,
  updated_at          TIMESTAMPTZ NOT NULL,
  deleted_at          TIMESTAMPTZ,         --　削除日時
  visibility          SMALLINT NOT NULL DEFAULT 1 CHECK (visibility IN (0,1,2)) -- 公開範囲0:private,1:unlisted,2:public
);

-- トリップの写真
CREATE TABLE route_images (
  id           BIGSERIAL PRIMARY KEY,
  ulid  CHAR(26) NOT NULL UNIQUE,
  route_id      BIGINT NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
  s3_key       TEXT NOT NULL,            -- S3等の保存先パス
  width        INTEGER,                             -- 画像の幅
  height       INTEGER,                             -- 画像の高さ
  size         BIGINT,                          -- ファイルサイズ（バイト）
  type         TEXT NOT NULL,                   -- jpg/png等
  visibility   SMALLINT NOT NULL DEFAULT 1 CHECK (visibility IN (0,1,2)),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (s3_key)
);

-- キューシート
CREATE TABLE  course_point(
  id            BIGSERIAL PRIMARY KEY,
  ulid  CHAR(26) NOT NULL UNIQUE,
  route_id      BIGINT REFERENCES routes(id) ON DELETE CASCADE,
  step_order    INT NOT NULL,          -- 0..n（ルート全体の通し順）
  seg_dist_m    DOUBLE PRECISION,      -- 直前のポイントからこのポイントまでの区間距離(m)
  cum_dist_m    DOUBLE PRECISION,      -- 直前のポイントからこのポイントまでの区間距離(m)
  duration      DOUBLE PRECISION,      -- 直前からこのポイyントまでの所要時間(s)
  instruction   TEXT,
  road_name     TEXT,
  maneuver_type TEXT,                  -- 'turn','depart','arrive'等
  modifier      TEXT,                  -- 'left','right','slight_left'等
  location      geometry(Point, 4326), 
  bearing_before INT,
  bearing_after  INT
);

-- 活動
CREATE TABLE trips (
  id                     BIGSERIAL PRIMARY KEY,                   -- 例: 342859653
  ulid  CHAR(26) NOT NULL UNIQUE,
  user_id                BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  -- 表示/識別
  name                   TEXT NOT NULL DEFAULT '',
  description            TEXT NOT NULL DEFAULT '',
  visibility             SMALLINT NOT NULL DEFAULT 0 CHECK (visibility IN (0,1,2)),
  highlighted_photo_id   BIGINT NOT NULL DEFAULT 0,
  
  -- 位置情報（PostGIS を使うなら Point/Polygon で保持）
  path_geom              geometry(LineString, 4326),
  first_point            geometry(Point, 4326),                -- ST_SetSRID(ST_MakePoint(lng,lat),4326)
  last_point             geometry(Point, 4326),
  bbox_geom              geometry(Polygon, 4326),              -- NE/SW から生成して保存も可

  -- 計測系（距離[m]、時間[秒]、速度[m/s] など）
  distance               DOUBLE PRECISION CHECK (distance IS NULL OR distance >= 0),                     
  duration               INTEGER CHECK (duration IS NULL OR duration >= 0),                              
  moving_time            INTEGER CHECK (moving_time IS NULL OR moving_time >= 0),                              
  elevation_gain         DOUBLE PRECISION CHECK (elevation_gain IS NULL OR elevation_gain >= 0),
  elevation_loss         DOUBLE PRECISION CHECK (elevation_loss IS NULL OR elevation_loss >= 0),
  avg_speed              DOUBLE PRECISION CHECK (avg_speed IS NULL OR avg_speed >= 0),
  max_speed              DOUBLE PRECISION CHECK (max_speed IS NULL OR max_speed >= 0),

  -- センサー/パワー関連
  avg_cad                DOUBLE PRECISION,    -- 平均ケイデンス[rpm]
  max_cad                DOUBLE PRECISION,
  min_cad                DOUBLE PRECISION,
  max_hr                 INTEGER,             -- 最大心拍[bpm]
  min_hr                 INTEGER,
  avg_watts              DOUBLE PRECISION,
  max_watts              DOUBLE PRECISION,
  min_watts              DOUBLE PRECISION,
  avg_watts_estimated    BOOLEAN,      -- センサー実測か速度・勾配からの推定値か区別するフラグ
  avg_power_estimated    DOUBLE PRECISION, -- 別のロジックで計算した推定平均パワーを入れるためのフィールド(null可）
  calories               DOUBLE PRECISION,  -- 消費カロリー[kcal]

  -- 種別/状態
  is_gps                 BOOLEAN NOT NULL DEFAULT FALSE,
  is_stationary          BOOLEAN NOT NULL DEFAULT FALSE,   -- 固定ローラー / スマートトレーナーなど「移動していない室内トレーニング」フラグ
  processed              BOOLEAN NOT NULL DEFAULT FALSE,

  -- 時刻/タイムゾーン
  created_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at             TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at             TIMESTAMPTZ,
  departed_at            TIMESTAMPTZ,
  time_zone              TEXT,                                  -- 例: 'Asia/Tokyo'
  utc_offset             INTEGER,                                -- 秒（例: +9h = 32400）

  -- アクティビティ
  activity_type_id       INTEGER NOT NULL DEFAULT 0,
  
  -- 付帯情報
  pace                   DOUBLE PRECISION,
  moving_pace            DOUBLE PRECISION
);

-- トリップの写真
CREATE TABLE trip_images (
  id           BIGSERIAL PRIMARY KEY,
  ulid  CHAR(26) NOT NULL UNIQUE,
  trip_id      BIGINT NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
  s3_key       TEXT NOT NULL,            -- S3等の保存先パス
  width        INTEGER,                             -- 画像の幅
  height       INTEGER,                             -- 画像の高さ
  size         BIGINT,                          -- ファイルサイズ（バイト）
  type         TEXT NOT NULL,                   -- jpg/png等
  visibility   SMALLINT NOT NULL DEFAULT 1 CHECK (visibility IN (0,1,2)),
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (s3_key)
);


CREATE TABLE route_likes (
  id           BIGSERIAL PRIMARY KEY,
  ulid  CHAR(26) NOT NULL UNIQUE,
  user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  route_id    BIGINT NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  UNIQUE (user_id, route_id)
);


CREATE TABLE route_comments (
  id           BIGSERIAL PRIMARY KEY,
  ulid  CHAR(26) NOT NULL UNIQUE,
  user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  route_id    BIGINT NOT NULL,
  parent_id    BIGINT,                            -- 返信ツリー（同テーブル参照）
  content      TEXT NOT NULL,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at   TIMESTAMPTZ,
  CONSTRAINT comments_parent_fk FOREIGN KEY (parent_id) REFERENCES route_comments(id) ON DELETE SET NULL
);

-- ルートのブックマーク（保存）
CREATE TABLE route_saves (
  id         BIGSERIAL PRIMARY KEY,
  ulid       CHAR(26) NOT NULL UNIQUE,
  user_id    BIGINT NOT NULL REFERENCES users(id)   ON DELETE CASCADE,
  route_id   BIGINT NOT NULL REFERENCES routes(id)  ON DELETE CASCADE,
  pinned     BOOLEAN NOT NULL DEFAULT FALSE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,               -- ソフト削除（履歴/復活用）

  UNIQUE (user_id, route_id)            -- 同じルートの重複保存を防ぐ
);

