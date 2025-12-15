-- SQL script to create Breaker snapshot table in Iceberg
-- Generated from scripts/generated_iceberg_tables.sql

-- Create namespace
CREATE NAMESPACE IF NOT EXISTS ontology.grid;

-- Create Breaker snapshot table
CREATE TABLE IF NOT EXISTS ontology.grid.breaker_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    locked STRING COMMENT 'If true the switch is locked. The resulting switch state is a combination of locked and Switch.open attributes as followsullilockedtrue and Switch.opentrue. The resulting state is open and lockedlililockedfalse and Switch.opentrue. The resulting state is openlililockedfalse and Switch.openfalse. The resulting state is closed.liul',
    normal_open STRING COMMENT 'The attribute is used in cases when no Measurement for the status value is present. If the Switch has a status measurement the Discrete.normalValue is expected to match with the Switch.normalOpen.',
    open STRING COMMENT 'The attribute tells if the switch is considered open when used as input to topology processing.',
    retained STRING COMMENT 'Branch is retained in the topological solution.  The flow through retained switches will normally be calculated in power flow.',
    switch_on_count STRING COMMENT 'The switch on count since the switch was last reset or initialized.',
    switch_on_date STRING COMMENT 'The date and time when the switch was last switched on.',
    aggregate STRING COMMENT 'The aggregate flag provides an alternative way of representing an aggregated equivalent element. It is applicable in cases when the dedicated classes for equivalent equipment do not have all of the attributes necessary to represent the required level of detail.  In case the flag is set to true the single instance of equipment represents multiple pieces of equipment that have been modelled together as an aggregate equivalent obtained by a network reduction procedure. Examples would be power transformers or synchronous machines operating in parallel modelled as a single aggregate power transformer or aggregate synchronous machine.  The attribute is not used for EquivalentBranch EquivalentShunt and EquivalentInjection.',
    in_service STRING COMMENT 'Specifies the availability of the equipment. True means the equipment is available for topology processing which determines if the equipment is energized or not. False means that the equipment is treated by network applications as if it is not in the model.',
    network_analysis_enabled STRING COMMENT 'The equipment is enabled to participate in network analysis.  If unspecified the value is assumed to be true.',
    normally_in_service STRING COMMENT 'Specifies the availability of the equipment under normal operating conditions. True means the equipment is available for topology processing which determines if the equipment is energized or not. False means that the equipment is treated by network applications as if it is not in the model.',
    alias_name STRING COMMENT 'The aliasName is free text human readable name of the object alternative to IdentifiedObject.name. It may be non unique and may not correlate to a naming hierarchy.The attribute aliasName is retained because of backwards compatibility between CIM relases. It is however recommended to replace aliasName with the Name class as aliasName is planned for retirement at a future time.',
    description STRING COMMENT 'The description is a free human readable text describing or naming the object. It may be non unique and may not correlate to a naming hierarchy.',
    m_rid STRING COMMENT 'Master resource identifier issued by a model authority. The mRID is unique within an exchange context. Global uniqueness is easily achieved by using a UUID as specified in RFC 4122 for the mRID. The use of UUID is strongly recommended.For CIMXML data files in RDF syntax conforming to IEC 61970552 the mRID is mapped to rdfID or rdfabout attributes that identify CIM object elements.',
    name STRING COMMENT 'The name is any free human readable and possibly non unique text naming the object.',
valid_from TIMESTAMP,
valid_to TIMESTAMP,
op_type STRING,
ingestion_ts TIMESTAMP
)
USING ICEBERG
PARTITIONED BY (actor_id)
COMMENT 'Snapshot table for Breaker Actor state persistence (includes inherited properties from: ProtectedSwitch)'

-- Verify table creation (execute separately)
-- SHOW TABLES IN ontology.grid LIKE 'breaker*';
