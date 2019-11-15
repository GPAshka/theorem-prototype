CREATE SCHEMA IF NOT EXISTS device;
ALTER SCHEMA device OWNER TO "Theorem";

CREATE TABLE IF NOT EXISTS device."Devices"
(
    "SerialNumber" character varying COLLATE pg_catalog."default" NOT NULL,
    "RegistrationDate" timestamp without time zone NOT NULL,
    "FirmwareVersion" character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "Devices_pkey" PRIMARY KEY ("SerialNumber")
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE device."Devices"
    OWNER to "Theorem";


CREATE SEQUENCE IF NOT EXISTS device."SensorData_Id_seq"
    INCREMENT 1
    START 1
    MINVALUE 1
    MAXVALUE 9223372036854775807
    CACHE 1;

ALTER SEQUENCE device."SensorData_Id_seq"
    OWNER TO "Theorem";


CREATE TABLE IF NOT EXISTS device."SensorData"
(
    "Id" integer NOT NULL DEFAULT nextval('device."SensorData_Id_seq"'::regclass),
    "Date" timestamp without time zone NOT NULL,
    "Temperature" numeric NOT NULL,
    "AirHumidity" numeric NOT NULL,
    "CarbonMonoxide" numeric NOT NULL,
    "HealthStatus" character varying COLLATE pg_catalog."default" NOT NULL,
    "DeviceSerialNumber" character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT "SensorData_pkey" PRIMARY KEY ("Id"),
    CONSTRAINT "FK_SensorData_Device" FOREIGN KEY ("DeviceSerialNumber")
        REFERENCES device."Devices" ("SerialNumber") MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE device."SensorData"
    OWNER to "Theorem";
