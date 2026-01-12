-- Modify "routes" table
ALTER TABLE "public"."routes" DROP COLUMN "has_course_points", ALTER COLUMN "bbox" DROP EXPRESSION, ALTER COLUMN "first_point" DROP EXPRESSION, ALTER COLUMN "last_point" DROP EXPRESSION, ALTER COLUMN "created_at" SET DEFAULT now(), ALTER COLUMN "updated_at" SET DEFAULT now();
-- Create "course_points" table
CREATE TABLE "public"."course_points" (
  "id" uuid NOT NULL,
  "route_id" uuid NOT NULL,
  "step_order" integer NOT NULL,
  "seg_dist_m" double precision NULL,
  "cum_dist_m" double precision NULL,
  "duration" double precision NULL,
  "instruction" text NULL,
  "road_name" text NULL,
  "maneuver_type" text NULL,
  "modifier" text NULL,
  "location" public.geometry(Point,4326) NULL,
  "bearing_before" integer NULL,
  "bearing_after" integer NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "course_points_route_id_step_order_key" UNIQUE ("route_id", "step_order"),
  CONSTRAINT "course_points_route_id_fkey" FOREIGN KEY ("route_id") REFERENCES "public"."routes" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create "waypoints" table
CREATE TABLE "public"."waypoints" (
  "id" uuid NOT NULL,
  "route_id" uuid NOT NULL,
  "location" public.geometry(Point,4326) NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY ("id"),
  CONSTRAINT "waypoints_route_id_fkey" FOREIGN KEY ("route_id") REFERENCES "public"."routes" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Drop "course_point" table
DROP TABLE "public"."course_point";
