CREATE NAMESPACE IF NOT EXISTS ontology.grid

CREATE NAMESPACE IF NOT EXISTS ontology.assets

CREATE NAMESPACE IF NOT EXISTS ontology.metering

CREATE NAMESPACE IF NOT EXISTS ontology.actors

CREATE TABLE IF NOT EXISTS ontology.grid.acdc_converter_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    number_of_valves STRING COMMENT 'Number of valves in the converter. Used in loss calculations.',
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
COMMENT 'Snapshot table for ACDCConverter Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.ac_line_segment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for ACLineSegment Actor state persistence (includes inherited properties from: Conductor)'


CREATE TABLE IF NOT EXISTS ontology.grid.ac_line_segment_phase_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    sequence_number STRING COMMENT 'Number designation for this line segment phase. Each line segment phase within a line segment should have a unique sequence number. This is useful for unbalanced modelling to bind the mathematical model PhaseImpedanceData of PerLengthPhaseImpedance with the connectivity model this class and the physical model WirePosition without tight coupling.',
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
COMMENT 'Snapshot table for ACLineSegmentPhase Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.air_compressor_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    air_compressor_rating STRING COMMENT 'Rating of the CAES air compressor.',
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
COMMENT 'Snapshot table for AirCompressor Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.asynchronous_machine_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    converter_fed_drive STRING COMMENT 'Indicates whether the machine is a converter fed drive. Used for short circuit data exchange according to IEC 60909.',
    ia_ir_ratio STRING COMMENT 'Ratio of lockedrotor current to the rated current of the motor IaIr. Used for short circuit data exchange according to IEC 60909.',
    pole_pair_number STRING COMMENT 'Number of pole pairs of stator. Used for short circuit data exchange according to IEC 60909.',
    reversible STRING COMMENT 'Indicates for converter drive motors if the power can be reversible. Used for short circuit data exchange according to IEC 60909.',
    rx_locked_rotor_ratio STRING COMMENT 'Locked rotor ratio RX. Used for short circuit data exchange according to IEC 60909.',
    rated_power_factor STRING COMMENT 'Power factor nameplate data. It is primarily used for short circuit data exchange according to IEC 60909. The attribute cannot be a negative value.',
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for AsynchronousMachine Actor state persistence (includes inherited properties from: RotatingMachine)'


CREATE TABLE IF NOT EXISTS ontology.grid.auxiliary_equipment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for AuxiliaryEquipment Actor state persistence (includes inherited properties from: Equipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.bwr_steam_supply_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    integral_gain STRING COMMENT 'Integral gain.',
    pressure_setpoint_ga STRING COMMENT 'Pressure setpoint gain adjuster.',
    proportional_gain STRING COMMENT 'Proportional gain.',
    rod_pattern_constant STRING COMMENT 'Constant associated with rod pattern.',
    steam_supply_rating STRING COMMENT 'Rating of steam supply.',
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
COMMENT 'Snapshot table for BWRSteamSupply Actor state persistence (includes inherited properties from: SteamSupply)'


CREATE TABLE IF NOT EXISTS ontology.grid.battery_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for BatteryUnit Actor state persistence (includes inherited properties from: PowerElectronicsUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.bay_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    bay_energy_meas_flag STRING COMMENT 'Indicates the presenceabsence of energy measurements.',
    bay_power_meas_flag STRING COMMENT 'Indicates the presenceabsence of activereactive power measurements.',
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
COMMENT 'Snapshot table for Bay Actor state persistence (includes inherited properties from: EquipmentContainer)'


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


CREATE TABLE IF NOT EXISTS ontology.grid.busbar_section_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for BusbarSection Actor state persistence (includes inherited properties from: Connector)'


CREATE TABLE IF NOT EXISTS ontology.grid.caes_plant_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for CAESPlant Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.circuit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Circuit Actor state persistence (includes inherited properties from: Line)'


CREATE TABLE IF NOT EXISTS ontology.grid.clamp_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    aggregate STRING COMMENT 'The aggregate flag provides an alternative way of representing an aggregated equivalent element. It is applicable in cases when the dedicated classes for equivalent equipment do not have all of the attributes necessary to represent the required level of detail.  In case the flag is set to true the single instance of equipment represents multiple pieces of equipment that have been modelled together as an aggregate equivalent obtained by a network reduction procedure. Examples would be power transformers or synchronous machines operating in parallel modelled as a single aggregate power transformer or aggregate synchronous machine.  The attribute is not used for EquivalentBranch EquivalentShunt and EquivalentInjection.',
    in_service STRING COMMENT 'Specifies the availability of the equipment. True means the equipment is available for topology processing which determines if the equipment is energized or not. False means that the equipment is treated by network applications as if it is not in the model.',
    network_analysis_enabled STRING COMMENT 'The equipment is enabled to participate in network analysis.  If unspecified the value is assumed to be true.',
    normally_in_service STRING COMMENT 'Specifies the availability of the equipment under normal operating conditions. True means the equipment is available for topology processing which determines if the equipment is energized or not. False means that the equipment is treated by network applications as if it is not in the model.',
    alias_name STRING COMMENT 'The aliasName is free text human readable name of the object alternative to IdentifiedObject.name. It may be non unique and may not correlate to a naming hierarchy.The attribute aliasName is retained because of backwards compatibility between CIM relases. It is however recommended to replace aliasName with the Name class as aliasName is planned for retirement at a future time.',
    description STRING COMMENT 'The description is a free human readable text describing or naming the object. It may be non unique and may not correlate to a naming hierarchy.',
    m_rid STRING COMMENT 'Master resource identifier issued by a model authority. The mRID is unique within an exchange context. Global uniqueness is easily achieved by using a UUID as specified in RFC 4122 for the mRID. The use of UUID is strongly recommended.For CIMXML data files in RDF syntax conforming to IEC 61970552 the mRID is mapped to rdfID or rdfabout attributes that identify CIM object elements.',
    name STRING COMMENT 'The name is any free human readable and possibly non unique text naming the object.',
    executed_date_time STRING COMMENT 'Actual date and time of this switching step.',
    issued_date_time STRING COMMENT 'Date and time when the crew was given the instruction to execute the action not applicable if the action is performed by operator remote control.',
    planned_date_time STRING COMMENT 'Planned date and time of this switching step.',
valid_from TIMESTAMP,
valid_to TIMESTAMP,
op_type STRING,
ingestion_ts TIMESTAMP
)
USING ICEBERG
PARTITIONED BY (actor_id)
COMMENT 'Snapshot table for Clamp Actor state persistence (includes inherited properties from: ConductingEquipment, SwitchingAction)'


CREATE TABLE IF NOT EXISTS ontology.grid.cogeneration_plant_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    cogen_hp_sendout_rating STRING COMMENT 'The high pressure steam sendout.',
    cogen_hp_steam_rating STRING COMMENT 'The high pressure steam rating.',
    cogen_lp_sendout_rating STRING COMMENT 'The low pressure steam sendout.',
    cogen_lp_steam_rating STRING COMMENT 'The low pressure steam rating.',
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
COMMENT 'Snapshot table for CogenerationPlant Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.combined_cycle_configuration_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    shutdown_flag STRING COMMENT 'Whether Combined Cycle Plant can be shutdown in this Configuration?',
    startup_flag STRING COMMENT 'Whether Combined Cycle Plant can be started in this Logical Configuration?',
    primary_configuration STRING COMMENT 'Whether this CombinedCycleConfiguration is the primary configuration in the associated Logical configuration?',
    cold_start_time STRING COMMENT 'Cold start time.',
    combined_cycle_operating_mode STRING COMMENT 'Combined Cycle operating mode.',
    commericial_operation_date STRING COMMENT '',
    hot_int_time STRING COMMENT 'Hottointermediate time Seasonal',
    hot_start_time STRING COMMENT 'Hot start time.',
    int_cold_time STRING COMMENT 'Intermediatetocold time Seasonal',
    int_start_time STRING COMMENT 'Intermediate start time.',
    max_shutdown_time STRING COMMENT 'Maximum time this device can be shut down.',
    max_start_ups_per_day STRING COMMENT 'maximum start ups per day',
    max_weekly_energy STRING COMMENT 'Maximum weekly Energy Seasonal',
    max_weekly_starts STRING COMMENT 'Maximum weekly starts seasonal parameter',
    pump_min_down_time STRING COMMENT 'The minimum down time for the pump in a pump storage unit.',
    pump_min_up_time STRING COMMENT 'The minimum up time aspect for the pump in a pump storage unit',
    pump_shutdown_cost STRING COMMENT 'The cost to shutdown a pump during the pump aspect of a pump storage unit.',
    pump_shutdown_time STRING COMMENT 'The shutdown time minutes of the pump aspect of a pump storage unit.',
    pumping_factor STRING COMMENT 'Pumping factor for pump storage units conversion factor between generating and pumping.',
    resource_sub_type STRING COMMENT 'Unit sub type used by Settlements or scheduling application. Application use of the unit sub type may define the necessary types as applicable.',
    river_system STRING COMMENT 'River System the Resource is tied to.',
    commercial_op_date STRING COMMENT 'Resource Commercial Operation Date. ',
    dispatchable STRING COMMENT 'Dispatchable indicates whether the resource is dispatchable. This implies that the resource intends to submit Energy bidsoffers or Ancillary Services bidsoffers or selfprovided schedules.',
    last_modified STRING COMMENT 'Indication of the last time this item was modifiedversioned.',
    max_base_self_sched_qty STRING COMMENT 'Maximum base self schedule quantity.',
    max_on_time STRING COMMENT 'Maximum on time after start up.',
    min_off_time STRING COMMENT 'Minimum off time after shut down.',
    min_on_time STRING COMMENT 'Minimum on time after start up.',
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
COMMENT 'Snapshot table for CombinedCycleConfiguration Actor state persistence (includes inherited properties from: RegisteredGenerator)'


CREATE TABLE IF NOT EXISTS ontology.grid.combined_cycle_plant_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for CombinedCyclePlant Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.combustion_turbine_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    heat_recovery_flag STRING COMMENT 'Flag that is set to true if the combustion turbine is associated with a heat recovery boiler.',
    prime_mover_rating STRING COMMENT 'Rating of prime mover.',
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
COMMENT 'Snapshot table for CombustionTurbine Actor state persistence (includes inherited properties from: PrimeMover)'


CREATE TABLE IF NOT EXISTS ontology.grid.communication_link_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for CommunicationLink Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.composite_switch_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    composite_switch_type STRING COMMENT 'An alphanumeric code that can be used as a reference to extra information such as the description of the interlocking scheme if any.',
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
COMMENT 'Snapshot table for CompositeSwitch Actor state persistence (includes inherited properties from: Equipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.conducting_equipment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for ConductingEquipment Actor state persistence (includes inherited properties from: Equipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.conductor_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Conductor Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.conform_load_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    customer_count STRING COMMENT 'Number of individual customers represented by this demand.',
    grounded STRING COMMENT 'Used for Yn and Zn connections. True if the neutral is solidly grounded.',
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
COMMENT 'Snapshot table for ConformLoad Actor state persistence (includes inherited properties from: EnergyConsumer)'


CREATE TABLE IF NOT EXISTS ontology.grid.connectivity_node_container_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for ConnectivityNodeContainer Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.connector_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Connector Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.control_area_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for ControlArea Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.cs_converter_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    number_of_valves STRING COMMENT 'Number of valves in the converter. Used in loss calculations.',
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
COMMENT 'Snapshot table for CsConverter Actor state persistence (includes inherited properties from: ACDCConverter)'


CREATE TABLE IF NOT EXISTS ontology.grid.current_relay_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    inverse_time_flag STRING COMMENT 'Set true if the current relay has inverse time characteristic.',
    high_limit STRING COMMENT 'The maximum allowable value.',
    low_limit STRING COMMENT 'The minimum allowable value.',
    power_direction_flag STRING COMMENT 'Direction same as positive active power flow value.',
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
COMMENT 'Snapshot table for CurrentRelay Actor state persistence (includes inherited properties from: ProtectionEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.current_transformer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    accuracy_class STRING COMMENT 'CT accuracy classification.',
    ct_class STRING COMMENT 'CT classification i.e. class 10P.',
    usage STRING COMMENT 'Intended usage of the CT i.e. metering protection.',
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
COMMENT 'Snapshot table for CurrentTransformer Actor state persistence (includes inherited properties from: Sensor)'


CREATE TABLE IF NOT EXISTS ontology.grid.cut_snapshots (
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
COMMENT 'Snapshot table for Cut Actor state persistence (includes inherited properties from: Switch)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_breaker_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCBreaker Actor state persistence (includes inherited properties from: DCSwitch)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_busbar_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCBusbar Actor state persistence (includes inherited properties from: DCConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_chopper_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCChopper Actor state persistence (includes inherited properties from: DCConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_conducting_equipment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCConductingEquipment Actor state persistence (includes inherited properties from: Equipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_converter_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCConverterUnit Actor state persistence (includes inherited properties from: DCEquipmentContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_disconnector_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCDisconnector Actor state persistence (includes inherited properties from: DCSwitch)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_equipment_container_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCEquipmentContainer Actor state persistence (includes inherited properties from: EquipmentContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_ground_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCGround Actor state persistence (includes inherited properties from: DCConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_line_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCLine Actor state persistence (includes inherited properties from: DCEquipmentContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_line_segment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCLineSegment Actor state persistence (includes inherited properties from: DCConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_series_device_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCSeriesDevice Actor state persistence (includes inherited properties from: DCConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_shunt_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCShunt Actor state persistence (includes inherited properties from: DCConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.dc_switch_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for DCSwitch Actor state persistence (includes inherited properties from: DCConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.disconnecting_circuit_breaker_snapshots (
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
COMMENT 'Snapshot table for DisconnectingCircuitBreaker Actor state persistence (includes inherited properties from: Breaker)'


CREATE TABLE IF NOT EXISTS ontology.grid.disconnector_snapshots (
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
COMMENT 'Snapshot table for Disconnector Actor state persistence (includes inherited properties from: Switch)'


CREATE TABLE IF NOT EXISTS ontology.grid.drum_boiler_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    drum_boiler_rating STRING COMMENT 'Rating of drum boiler in steam units.',
    control_error_bias_p STRING COMMENT 'Active power error bias ratio.',
    control_ic STRING COMMENT 'Integral constant.',
    control_pc STRING COMMENT 'Proportional constant.',
    control_peb STRING COMMENT 'Pressure error bias ratio.',
    control_tc STRING COMMENT 'Time constant.',
    feed_water_ig STRING COMMENT 'Feedwater integral gain ratio.',
    feed_water_pg STRING COMMENT 'Feedwater proportional gain ratio.',
    max_error_rate_p STRING COMMENT 'Active power maximum error rate limit.',
    min_error_rate_p STRING COMMENT 'Active power minimum error rate limit.',
    pressure_ctrl_dg STRING COMMENT 'Pressure control derivative gain ratio.',
    pressure_ctrl_ig STRING COMMENT 'Pressure control integral gain ratio.',
    pressure_ctrl_pg STRING COMMENT 'Pressure control proportional gain ratio.',
    pressure_feedback STRING COMMENT 'Pressure feedback indicator.',
    super_heater1_capacity STRING COMMENT 'Drumprimary superheater capacity.',
    super_heater2_capacity STRING COMMENT 'Secondary superheater capacity.',
    super_heater_pipe_pd STRING COMMENT 'Superheater pipe pressure drop constant.',
    steam_supply_rating STRING COMMENT 'Rating of steam supply.',
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
COMMENT 'Snapshot table for DrumBoiler Actor state persistence (includes inherited properties from: FossilSteamSupply)'


CREATE TABLE IF NOT EXISTS ontology.grid.earth_fault_compensator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EarthFaultCompensator Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.energy_connection_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EnergyConnection Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.energy_consumer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    customer_count STRING COMMENT 'Number of individual customers represented by this demand.',
    grounded STRING COMMENT 'Used for Yn and Zn connections. True if the neutral is solidly grounded.',
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
COMMENT 'Snapshot table for EnergyConsumer Actor state persistence (includes inherited properties from: EnergyConnection)'


CREATE TABLE IF NOT EXISTS ontology.grid.energy_consumer_phase_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EnergyConsumerPhase Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.energy_group_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    is_slack STRING COMMENT '',
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
COMMENT 'Snapshot table for EnergyGroup Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.energy_source_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EnergySource Actor state persistence (includes inherited properties from: EnergyConnection)'


CREATE TABLE IF NOT EXISTS ontology.grid.energy_source_phase_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EnergySourcePhase Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.equipment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Equipment Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.equipment_container_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EquipmentContainer Actor state persistence (includes inherited properties from: ConnectivityNodeContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.equivalent_branch_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EquivalentBranch Actor state persistence (includes inherited properties from: EquivalentEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.equivalent_equipment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EquivalentEquipment Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.equivalent_injection_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    regulation_capability STRING COMMENT 'Specifies whether or not the EquivalentInjection has the capability to regulate the local voltage. If true the EquivalentInjection can regulate. If false the EquivalentInjection cannot regulate. ReactiveCapabilityCurve can only be associated with EquivalentInjection  if the flag is true.',
    regulation_status STRING COMMENT 'Specifies the regulation status of the EquivalentInjection.  True is regulating.  False is not regulating.',
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
COMMENT 'Snapshot table for EquivalentInjection Actor state persistence (includes inherited properties from: EquivalentEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.equivalent_network_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EquivalentNetwork Actor state persistence (includes inherited properties from: ConnectivityNodeContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.equivalent_shunt_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for EquivalentShunt Actor state persistence (includes inherited properties from: EquivalentEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.external_network_injection_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    ik_second STRING COMMENT 'Indicates whether initial symmetrical shortcircuit current and power have been calculated according to IEC Ik.  Used only if short circuit calculations are done according to superposition method.',
    max_r0_to_x0_ratio STRING COMMENT 'Maximum ratio of zero sequence resistance of Network Feeder to its zero sequence reactance R0X0 max. Used for short circuit data exchange according to IEC 60909.',
    max_r1_to_x1_ratio STRING COMMENT 'Maximum ratio of positive sequence resistance of Network Feeder to its positive sequence reactance R1X1 max. Used for short circuit data exchange according to IEC 60909.',
    max_z0_to_z1_ratio STRING COMMENT 'Maximum ratio of zero sequence impedance to its positive sequence impedance Z0Z1 max. Used for short circuit data exchange according to IEC 60909.',
    min_r0_to_x0_ratio STRING COMMENT 'Indicates whether initial symmetrical shortcircuit current and power have been calculated according to IEC Ik. Used for short circuit data exchange according to IEC 6090.',
    min_r1_to_x1_ratio STRING COMMENT 'Minimum ratio of positive sequence resistance of Network Feeder to its positive sequence reactance R1X1 min. Used for short circuit data exchange according to IEC 60909.',
    min_z0_to_z1_ratio STRING COMMENT 'Minimum ratio of zero sequence impedance to its positive sequence impedance Z0Z1 min. Used for short circuit data exchange according to IEC 60909.',
    reference_priority STRING COMMENT 'Priority of unit for use as powerflow voltage phase angle reference bus selection. 0  don t care default 1  highest priority. 2 is less than 1 and so on.',
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for ExternalNetworkInjection Actor state persistence (includes inherited properties from: RegulatingCondEq)'


CREATE TABLE IF NOT EXISTS ontology.grid.fault_indicator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for FaultIndicator Actor state persistence (includes inherited properties from: AuxiliaryEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.feeder_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Feeder Actor state persistence (includes inherited properties from: EquipmentContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.flowgate_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Flowgate Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.fossil_steam_supply_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_error_bias_p STRING COMMENT 'Active power error bias ratio.',
    control_ic STRING COMMENT 'Integral constant.',
    control_pc STRING COMMENT 'Proportional constant.',
    control_peb STRING COMMENT 'Pressure error bias ratio.',
    control_tc STRING COMMENT 'Time constant.',
    feed_water_ig STRING COMMENT 'Feedwater integral gain ratio.',
    feed_water_pg STRING COMMENT 'Feedwater proportional gain ratio.',
    max_error_rate_p STRING COMMENT 'Active power maximum error rate limit.',
    min_error_rate_p STRING COMMENT 'Active power minimum error rate limit.',
    pressure_ctrl_dg STRING COMMENT 'Pressure control derivative gain ratio.',
    pressure_ctrl_ig STRING COMMENT 'Pressure control integral gain ratio.',
    pressure_ctrl_pg STRING COMMENT 'Pressure control proportional gain ratio.',
    pressure_feedback STRING COMMENT 'Pressure feedback indicator.',
    super_heater1_capacity STRING COMMENT 'Drumprimary superheater capacity.',
    super_heater2_capacity STRING COMMENT 'Secondary superheater capacity.',
    super_heater_pipe_pd STRING COMMENT 'Superheater pipe pressure drop constant.',
    steam_supply_rating STRING COMMENT 'Rating of steam supply.',
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
COMMENT 'Snapshot table for FossilSteamSupply Actor state persistence (includes inherited properties from: SteamSupply)'


CREATE TABLE IF NOT EXISTS ontology.grid.frequency_converter_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for FrequencyConverter Actor state persistence (includes inherited properties from: RegulatingCondEq)'


CREATE TABLE IF NOT EXISTS ontology.grid.fuse_snapshots (
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
COMMENT 'Snapshot table for Fuse Actor state persistence (includes inherited properties from: Switch)'


CREATE TABLE IF NOT EXISTS ontology.grid.generating_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    long_pf STRING COMMENT 'Generating unit long term economic participation factor.',
    normal_pf STRING COMMENT 'Generating unit economic participation factor.  The sum of the participation factors across generating units does not have to sum to one.  It is used for representing distributed slack participation factor. The attribute shall be a positive value or zero.',
    penalty_factor STRING COMMENT 'Defined as 1   1  Incremental Transmission Loss with the Incremental Transmission Loss expressed as a plus or minus value. The typical range of penalty factors is 0.9 to 1.1.',
    short_pf STRING COMMENT 'Generating unit short term economic participation factor.',
    tie_line_pf STRING COMMENT 'Generating unit economic participation factor.',
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
COMMENT 'Snapshot table for GeneratingUnit Actor state persistence (includes inherited properties from: Equipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.ground_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Ground Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.ground_disconnector_snapshots (
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
COMMENT 'Snapshot table for GroundDisconnector Actor state persistence (includes inherited properties from: Switch)'


CREATE TABLE IF NOT EXISTS ontology.grid.grounding_impedance_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for GroundingImpedance Actor state persistence (includes inherited properties from: EarthFaultCompensator)'


CREATE TABLE IF NOT EXISTS ontology.grid.heat_recovery_boiler_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    steam_supply_rating2 STRING COMMENT 'The steam supply rating in kilopounds per hour if dual pressure boiler.',
    control_error_bias_p STRING COMMENT 'Active power error bias ratio.',
    control_ic STRING COMMENT 'Integral constant.',
    control_pc STRING COMMENT 'Proportional constant.',
    control_peb STRING COMMENT 'Pressure error bias ratio.',
    control_tc STRING COMMENT 'Time constant.',
    feed_water_ig STRING COMMENT 'Feedwater integral gain ratio.',
    feed_water_pg STRING COMMENT 'Feedwater proportional gain ratio.',
    max_error_rate_p STRING COMMENT 'Active power maximum error rate limit.',
    min_error_rate_p STRING COMMENT 'Active power minimum error rate limit.',
    pressure_ctrl_dg STRING COMMENT 'Pressure control derivative gain ratio.',
    pressure_ctrl_ig STRING COMMENT 'Pressure control integral gain ratio.',
    pressure_ctrl_pg STRING COMMENT 'Pressure control proportional gain ratio.',
    pressure_feedback STRING COMMENT 'Pressure feedback indicator.',
    super_heater1_capacity STRING COMMENT 'Drumprimary superheater capacity.',
    super_heater2_capacity STRING COMMENT 'Secondary superheater capacity.',
    super_heater_pipe_pd STRING COMMENT 'Superheater pipe pressure drop constant.',
    steam_supply_rating STRING COMMENT 'Rating of steam supply.',
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
COMMENT 'Snapshot table for HeatRecoveryBoiler Actor state persistence (includes inherited properties from: FossilSteamSupply)'


CREATE TABLE IF NOT EXISTS ontology.grid.host_control_area_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    frequency_bias_factor STRING COMMENT 'The control areas frequency bias factor in MW0.1 Hz for automatic generation control AGC',
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
COMMENT 'Snapshot table for HostControlArea Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.hydro_generating_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    long_pf STRING COMMENT 'Generating unit long term economic participation factor.',
    normal_pf STRING COMMENT 'Generating unit economic participation factor.  The sum of the participation factors across generating units does not have to sum to one.  It is used for representing distributed slack participation factor. The attribute shall be a positive value or zero.',
    penalty_factor STRING COMMENT 'Defined as 1   1  Incremental Transmission Loss with the Incremental Transmission Loss expressed as a plus or minus value. The typical range of penalty factors is 0.9 to 1.1.',
    short_pf STRING COMMENT 'Generating unit short term economic participation factor.',
    tie_line_pf STRING COMMENT 'Generating unit economic participation factor.',
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
COMMENT 'Snapshot table for HydroGeneratingUnit Actor state persistence (includes inherited properties from: GeneratingUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.hydro_power_plant_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    penstock_type STRING COMMENT 'Type and configuration of hydro plant penstocks.',
    surge_tank_code STRING COMMENT 'A code describing the type or absence of surge tank that is associated with the hydro power plant.',
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
COMMENT 'Snapshot table for HydroPowerPlant Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.hydro_pump_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for HydroPump Actor state persistence (includes inherited properties from: Equipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.hydro_turbine_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    gate_rate_limit STRING COMMENT 'Gate rate limit.',
    prime_mover_rating STRING COMMENT 'Rating of prime mover.',
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
COMMENT 'Snapshot table for HydroTurbine Actor state persistence (includes inherited properties from: PrimeMover)'


CREATE TABLE IF NOT EXISTS ontology.grid.ip_access_point_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    address STRING COMMENT 'Is the dotted decimal IP Address resolve the IP address. The format is controlled by the value of the addressType.',
    gateway STRING COMMENT 'Is the dotted decimal IPAddress of the first hop router.  Format is controlled by the addressType.',
    subnet STRING COMMENT 'This is the IP subnet mask which controls the local vs nonlocal routing.',
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
COMMENT 'Snapshot table for IPAccessPoint Actor state persistence (includes inherited properties from: CommunicationLink)'


CREATE TABLE IF NOT EXISTS ontology.grid.iso_upper_layer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    ae_invoke STRING COMMENT 'Is part of the Application Entity addressing as specified by ISO Addressing.',
    ae_qual STRING COMMENT 'Is the AE qualifier and represents further application level addressing information.',
    ap_invoke STRING COMMENT 'Is a further application level OSI addressing parameter.',
    ap_title STRING COMMENT 'Is a sequence of integer strings separated by ..  The value in conjunction with other application addressing attributes e.g. other APs are used to select a specific application e.g. the ICCP application entity per the OSI reference model.  The sequence and its values represent a namespace whose values are governed by ISOIEC 74983.',
    osi_psel STRING COMMENT 'Is the addressing selector for OSI presentation addressing.',
    osi_ssel STRING COMMENT 'Is the OSI session layer addressing information.',
    osi_tsel STRING COMMENT 'Is the OSI Transport Layer addressing information.',
    keep_alive_time STRING COMMENT 'Indicates the default interval at which TCP will check if the TCP connection is still valid.',
    port STRING COMMENT 'This value is only needed to be specified for called nodes e.g. those that respond to a TCP.Open request.This value specifies the TCP port to be used. Well known and registered ports are preferred and can be found at httpwww.iana.orgassignmentsservicenamesportnumbersservicenamesportnumbers.xhtmlFor IEC 608706 TASE.2 e.g. ICCP and IEC 61850 the value used shall be 102 for nonTLS protected exchanges. The value shall be 3782 for TLS transported ICCP and 61850 exchanges.',
    address STRING COMMENT 'Is the dotted decimal IP Address resolve the IP address. The format is controlled by the value of the addressType.',
    gateway STRING COMMENT 'Is the dotted decimal IPAddress of the first hop router.  Format is controlled by the addressType.',
    subnet STRING COMMENT 'This is the IP subnet mask which controls the local vs nonlocal routing.',
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
COMMENT 'Snapshot table for ISOUpperLayer Actor state persistence (includes inherited properties from: TCPAccessPoint)'


CREATE TABLE IF NOT EXISTS ontology.grid.jumper_snapshots (
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
COMMENT 'Snapshot table for Jumper Actor state persistence (includes inherited properties from: Switch)'


CREATE TABLE IF NOT EXISTS ontology.grid.junction_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Junction Actor state persistence (includes inherited properties from: Connector)'


CREATE TABLE IF NOT EXISTS ontology.grid.line_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Line Actor state persistence (includes inherited properties from: EquipmentContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.linear_shunt_compensator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    grounded STRING COMMENT 'Used for Yn and Zn connections. True if the neutral is solidly grounded.',
    maximum_sections STRING COMMENT 'The maximum number of sections that may be switched in. ',
    normal_sections STRING COMMENT 'The normal number of sections switched in. The value shall be between zero and ShuntCompensator.maximumSections.',
    sections STRING COMMENT 'Shunt compensator sections in use. Starting value for steady state solution. The attribute shall be a positive value or zero. Non integer values are allowed to support continuous variables. The reasons for continuous value are to support study cases where no discrete shunt compensators has yet been designed a solutions where a narrow voltage band force the sections to oscillate or accommodate for a continuous solution as input. For LinearShuntConpensator the value shall be between zero and ShuntCompensator.maximumSections. At value zero the shunt compensator conductance and admittance is zero. Linear interpolation of conductance and admittance between the previous and next integer section is applied in case of noninteger values.For NonlinearShuntCompensators shall only be set to one of the NonlinearShuntCompenstorPoint.sectionNumber. There is no interpolation between NonlinearShuntCompenstorPoints.',
    switch_on_count STRING COMMENT 'The switch on count since the capacitor count was last reset or initialized.',
    switch_on_date STRING COMMENT 'The date and time when the capacitor bank was last switched on.',
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for LinearShuntCompensator Actor state persistence (includes inherited properties from: ShuntCompensator)'


CREATE TABLE IF NOT EXISTS ontology.grid.linear_shunt_compensator_phase_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    maximum_sections STRING COMMENT 'The maximum number of sections that may be switched in for this phase. ',
    normal_sections STRING COMMENT 'For the capacitor phase the normal number of sections switched in. The value shall be between zero and ShuntCompensatorPhase.maximumSections.',
    sections STRING COMMENT 'Shunt compensator sections in use. Starting value for steady state solution. The attribute shall be a positive value or zero. Non integer values are allowed to support continuous variables. The reasons for continuous value are to support study cases where no discrete shunt compensators has yet been designed a solutions where a narrow voltage band force the sections to oscillate or accommodate for a continuous solution as input.For LinearShuntConpensator the value shall be between zero and ShuntCompensatorPhase.maximumSections. At value zero the shunt compensator conductance and admittance is zero. Linear interpolation of conductance and admittance between the previous and next integer section is applied in case of noninteger values.For NonlinearShuntCompensators shall only be set to one of the NonlinearShuntCompenstorPhasePoint.sectionNumber. There is no interpolation between NonlinearShuntCompenstorPhasePoints.',
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
COMMENT 'Snapshot table for LinearShuntCompensatorPhase Actor state persistence (includes inherited properties from: ShuntCompensatorPhase)'


CREATE TABLE IF NOT EXISTS ontology.grid.load_break_switch_snapshots (
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
COMMENT 'Snapshot table for LoadBreakSwitch Actor state persistence (includes inherited properties from: ProtectedSwitch)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_ac_line_segment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for MktACLineSegment Actor state persistence (includes inherited properties from: ACLineSegment)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_combined_cycle_plant_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for MktCombinedCyclePlant Actor state persistence (includes inherited properties from: CombinedCyclePlant)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_conducting_equipment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for MktConductingEquipment Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_control_area_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for MktControlArea Actor state persistence (includes inherited properties from: ControlArea)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_generating_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    long_pf STRING COMMENT 'Generating unit long term economic participation factor.',
    normal_pf STRING COMMENT 'Generating unit economic participation factor.  The sum of the participation factors across generating units does not have to sum to one.  It is used for representing distributed slack participation factor. The attribute shall be a positive value or zero.',
    penalty_factor STRING COMMENT 'Defined as 1   1  Incremental Transmission Loss with the Incremental Transmission Loss expressed as a plus or minus value. The typical range of penalty factors is 0.9 to 1.1.',
    short_pf STRING COMMENT 'Generating unit short term economic participation factor.',
    tie_line_pf STRING COMMENT 'Generating unit economic participation factor.',
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
COMMENT 'Snapshot table for MktGeneratingUnit Actor state persistence (includes inherited properties from: GeneratingUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_line_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for MktLine Actor state persistence (includes inherited properties from: Line)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_power_transformer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    is_part_of_generator_unit STRING COMMENT 'Indicates whether the machine is part of a power station unit. Used for short circuit data exchange according to IEC 60909.  It has an impact on how the correction factors are calculated for transformers since the transformer is not necessarily part of a synchronous machine and generating unit. It is not always possible to derive this information from the model. This is why the attribute is necessary.',
    operational_values_considered STRING COMMENT 'It is used to define if the data other attributes related to short circuit data exchange defines long term operational conditions or not. Used for short circuit data exchange according to IEC 60909.',
    vector_group STRING COMMENT 'Vector group of the transformer for protective relaying e.g. Dyn1. For unbalanced transformers this may not be simply determined from the constituent winding connections and phase angle displacements.The vectorGroup string consists of the following components in the order listed high voltage winding connection mid voltage winding connection for three winding transformers phase displacement clock number from 0 to 11  low voltage winding connection phase displacement clock number from 0 to 11.   The winding connections are D delta Y wye YN wye with neutral Z zigzag ZN zigzag with neutral A auto transformer. Upper case means the high voltage lower case mid or low. The high voltage winding always has clock position 0 and is not included in the vector group string.  Some examples YNy0 two winding wye to wye with no phase displacement YNd11 two winding wye to delta with 330 degrees phase displacement YNyn0d5 three winding transformer wye with neutral high voltage wye with neutral mid voltage and no phase displacement delta low voltage with 150 degrees displacement.Phase displacement is defined as the angular difference between the phasors representing the voltages between the neutral point real or imaginary and the corresponding terminals of two windings a positive sequence voltage system being applied to the highvoltage terminals following each other in alphabetical sequence if they are lettered or in numerical sequence if they are numbered the phasors are assumed to rotate in a counterclockwise sense.',
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
COMMENT 'Snapshot table for MktPowerTransformer Actor state persistence (includes inherited properties from: PowerTransformer)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_series_compensator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    varistor_present STRING COMMENT 'Describe if a metal oxide varistor mov for over voltage protection is configured in parallel with the series compensator. It is used for short circuit calculations.',
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
COMMENT 'Snapshot table for MktSeriesCompensator Actor state persistence (includes inherited properties from: SeriesCompensator)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_shunt_compensator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    grounded STRING COMMENT 'Used for Yn and Zn connections. True if the neutral is solidly grounded.',
    maximum_sections STRING COMMENT 'The maximum number of sections that may be switched in. ',
    normal_sections STRING COMMENT 'The normal number of sections switched in. The value shall be between zero and ShuntCompensator.maximumSections.',
    sections STRING COMMENT 'Shunt compensator sections in use. Starting value for steady state solution. The attribute shall be a positive value or zero. Non integer values are allowed to support continuous variables. The reasons for continuous value are to support study cases where no discrete shunt compensators has yet been designed a solutions where a narrow voltage band force the sections to oscillate or accommodate for a continuous solution as input. For LinearShuntConpensator the value shall be between zero and ShuntCompensator.maximumSections. At value zero the shunt compensator conductance and admittance is zero. Linear interpolation of conductance and admittance between the previous and next integer section is applied in case of noninteger values.For NonlinearShuntCompensators shall only be set to one of the NonlinearShuntCompenstorPoint.sectionNumber. There is no interpolation between NonlinearShuntCompenstorPoints.',
    switch_on_count STRING COMMENT 'The switch on count since the capacitor count was last reset or initialized.',
    switch_on_date STRING COMMENT 'The date and time when the capacitor bank was last switched on.',
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for MktShuntCompensator Actor state persistence (includes inherited properties from: ShuntCompensator)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_switch_snapshots (
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
COMMENT 'Snapshot table for MktSwitch Actor state persistence (includes inherited properties from: Switch)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_tap_changer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
    high_step STRING COMMENT 'Highest possible tap step position advance from neutral.The attribute shall be greater than lowStep.',
    low_step STRING COMMENT 'Lowest possible tap step position retard from neutral.',
    ltc_flag STRING COMMENT 'Specifies whether or not a TapChanger has load tap changing capabilities.',
    neutral_step STRING COMMENT 'The neutral tap step position for this winding.The attribute shall be equal to or greater than lowStep and equal or less than highStep.It is the step position where the voltage is neutralU when the other terminals of the transformer are at the ratedU.  If there are other tap changers on the transformer those taps are kept constant at their neutralStep.',
    normal_step STRING COMMENT 'The tap step position used in normal network operation for this winding. For a Fixed tap changer indicates the current physical tap setting.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
    step STRING COMMENT 'Tap changer position.Starting step for a steady state solution. Non integer values are allowed to support continuous tap variables. The reasons for continuous value are to support study cases where no discrete tap changer has yet been designed a solution where a narrow voltage band forces the tap step to oscillate or to accommodate for a continuous solution as input.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
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
COMMENT 'Snapshot table for MktTapChanger Actor state persistence (includes inherited properties from: TapChanger)'


CREATE TABLE IF NOT EXISTS ontology.grid.mkt_thermal_generating_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    long_pf STRING COMMENT 'Generating unit long term economic participation factor.',
    normal_pf STRING COMMENT 'Generating unit economic participation factor.  The sum of the participation factors across generating units does not have to sum to one.  It is used for representing distributed slack participation factor. The attribute shall be a positive value or zero.',
    penalty_factor STRING COMMENT 'Defined as 1   1  Incremental Transmission Loss with the Incremental Transmission Loss expressed as a plus or minus value. The typical range of penalty factors is 0.9 to 1.1.',
    short_pf STRING COMMENT 'Generating unit short term economic participation factor.',
    tie_line_pf STRING COMMENT 'Generating unit economic participation factor.',
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
COMMENT 'Snapshot table for MktThermalGeneratingUnit Actor state persistence (includes inherited properties from: ThermalGeneratingUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.non_conform_load_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    customer_count STRING COMMENT 'Number of individual customers represented by this demand.',
    grounded STRING COMMENT 'Used for Yn and Zn connections. True if the neutral is solidly grounded.',
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
COMMENT 'Snapshot table for NonConformLoad Actor state persistence (includes inherited properties from: EnergyConsumer)'


CREATE TABLE IF NOT EXISTS ontology.grid.nonlinear_shunt_compensator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    grounded STRING COMMENT 'Used for Yn and Zn connections. True if the neutral is solidly grounded.',
    maximum_sections STRING COMMENT 'The maximum number of sections that may be switched in. ',
    normal_sections STRING COMMENT 'The normal number of sections switched in. The value shall be between zero and ShuntCompensator.maximumSections.',
    sections STRING COMMENT 'Shunt compensator sections in use. Starting value for steady state solution. The attribute shall be a positive value or zero. Non integer values are allowed to support continuous variables. The reasons for continuous value are to support study cases where no discrete shunt compensators has yet been designed a solutions where a narrow voltage band force the sections to oscillate or accommodate for a continuous solution as input. For LinearShuntConpensator the value shall be between zero and ShuntCompensator.maximumSections. At value zero the shunt compensator conductance and admittance is zero. Linear interpolation of conductance and admittance between the previous and next integer section is applied in case of noninteger values.For NonlinearShuntCompensators shall only be set to one of the NonlinearShuntCompenstorPoint.sectionNumber. There is no interpolation between NonlinearShuntCompenstorPoints.',
    switch_on_count STRING COMMENT 'The switch on count since the capacitor count was last reset or initialized.',
    switch_on_date STRING COMMENT 'The date and time when the capacitor bank was last switched on.',
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for NonlinearShuntCompensator Actor state persistence (includes inherited properties from: ShuntCompensator)'


CREATE TABLE IF NOT EXISTS ontology.grid.nonlinear_shunt_compensator_phase_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    maximum_sections STRING COMMENT 'The maximum number of sections that may be switched in for this phase. ',
    normal_sections STRING COMMENT 'For the capacitor phase the normal number of sections switched in. The value shall be between zero and ShuntCompensatorPhase.maximumSections.',
    sections STRING COMMENT 'Shunt compensator sections in use. Starting value for steady state solution. The attribute shall be a positive value or zero. Non integer values are allowed to support continuous variables. The reasons for continuous value are to support study cases where no discrete shunt compensators has yet been designed a solutions where a narrow voltage band force the sections to oscillate or accommodate for a continuous solution as input.For LinearShuntConpensator the value shall be between zero and ShuntCompensatorPhase.maximumSections. At value zero the shunt compensator conductance and admittance is zero. Linear interpolation of conductance and admittance between the previous and next integer section is applied in case of noninteger values.For NonlinearShuntCompensators shall only be set to one of the NonlinearShuntCompenstorPhasePoint.sectionNumber. There is no interpolation between NonlinearShuntCompenstorPhasePoints.',
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
COMMENT 'Snapshot table for NonlinearShuntCompensatorPhase Actor state persistence (includes inherited properties from: ShuntCompensatorPhase)'


CREATE TABLE IF NOT EXISTS ontology.grid.nuclear_generating_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    long_pf STRING COMMENT 'Generating unit long term economic participation factor.',
    normal_pf STRING COMMENT 'Generating unit economic participation factor.  The sum of the participation factors across generating units does not have to sum to one.  It is used for representing distributed slack participation factor. The attribute shall be a positive value or zero.',
    penalty_factor STRING COMMENT 'Defined as 1   1  Incremental Transmission Loss with the Incremental Transmission Loss expressed as a plus or minus value. The typical range of penalty factors is 0.9 to 1.1.',
    short_pf STRING COMMENT 'Generating unit short term economic participation factor.',
    tie_line_pf STRING COMMENT 'Generating unit economic participation factor.',
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
COMMENT 'Snapshot table for NuclearGeneratingUnit Actor state persistence (includes inherited properties from: GeneratingUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.pwr_steam_supply_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    steam_supply_rating STRING COMMENT 'Rating of steam supply.',
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
COMMENT 'Snapshot table for PWRSteamSupply Actor state persistence (includes inherited properties from: SteamSupply)'


CREATE TABLE IF NOT EXISTS ontology.grid.petersen_coil_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for PetersenCoil Actor state persistence (includes inherited properties from: EarthFaultCompensator)'


CREATE TABLE IF NOT EXISTS ontology.grid.phase_tap_changer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
    high_step STRING COMMENT 'Highest possible tap step position advance from neutral.The attribute shall be greater than lowStep.',
    low_step STRING COMMENT 'Lowest possible tap step position retard from neutral.',
    ltc_flag STRING COMMENT 'Specifies whether or not a TapChanger has load tap changing capabilities.',
    neutral_step STRING COMMENT 'The neutral tap step position for this winding.The attribute shall be equal to or greater than lowStep and equal or less than highStep.It is the step position where the voltage is neutralU when the other terminals of the transformer are at the ratedU.  If there are other tap changers on the transformer those taps are kept constant at their neutralStep.',
    normal_step STRING COMMENT 'The tap step position used in normal network operation for this winding. For a Fixed tap changer indicates the current physical tap setting.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
    step STRING COMMENT 'Tap changer position.Starting step for a steady state solution. Non integer values are allowed to support continuous tap variables. The reasons for continuous value are to support study cases where no discrete tap changer has yet been designed a solution where a narrow voltage band forces the tap step to oscillate or to accommodate for a continuous solution as input.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
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
COMMENT 'Snapshot table for PhaseTapChanger Actor state persistence (includes inherited properties from: TapChanger)'


CREATE TABLE IF NOT EXISTS ontology.grid.phase_tap_changer_asymmetrical_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
    high_step STRING COMMENT 'Highest possible tap step position advance from neutral.The attribute shall be greater than lowStep.',
    low_step STRING COMMENT 'Lowest possible tap step position retard from neutral.',
    ltc_flag STRING COMMENT 'Specifies whether or not a TapChanger has load tap changing capabilities.',
    neutral_step STRING COMMENT 'The neutral tap step position for this winding.The attribute shall be equal to or greater than lowStep and equal or less than highStep.It is the step position where the voltage is neutralU when the other terminals of the transformer are at the ratedU.  If there are other tap changers on the transformer those taps are kept constant at their neutralStep.',
    normal_step STRING COMMENT 'The tap step position used in normal network operation for this winding. For a Fixed tap changer indicates the current physical tap setting.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
    step STRING COMMENT 'Tap changer position.Starting step for a steady state solution. Non integer values are allowed to support continuous tap variables. The reasons for continuous value are to support study cases where no discrete tap changer has yet been designed a solution where a narrow voltage band forces the tap step to oscillate or to accommodate for a continuous solution as input.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
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
COMMENT 'Snapshot table for PhaseTapChangerAsymmetrical Actor state persistence (includes inherited properties from: PhaseTapChangerNonLinear)'


CREATE TABLE IF NOT EXISTS ontology.grid.phase_tap_changer_linear_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
    high_step STRING COMMENT 'Highest possible tap step position advance from neutral.The attribute shall be greater than lowStep.',
    low_step STRING COMMENT 'Lowest possible tap step position retard from neutral.',
    ltc_flag STRING COMMENT 'Specifies whether or not a TapChanger has load tap changing capabilities.',
    neutral_step STRING COMMENT 'The neutral tap step position for this winding.The attribute shall be equal to or greater than lowStep and equal or less than highStep.It is the step position where the voltage is neutralU when the other terminals of the transformer are at the ratedU.  If there are other tap changers on the transformer those taps are kept constant at their neutralStep.',
    normal_step STRING COMMENT 'The tap step position used in normal network operation for this winding. For a Fixed tap changer indicates the current physical tap setting.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
    step STRING COMMENT 'Tap changer position.Starting step for a steady state solution. Non integer values are allowed to support continuous tap variables. The reasons for continuous value are to support study cases where no discrete tap changer has yet been designed a solution where a narrow voltage band forces the tap step to oscillate or to accommodate for a continuous solution as input.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
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
COMMENT 'Snapshot table for PhaseTapChangerLinear Actor state persistence (includes inherited properties from: PhaseTapChanger)'


CREATE TABLE IF NOT EXISTS ontology.grid.phase_tap_changer_non_linear_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
    high_step STRING COMMENT 'Highest possible tap step position advance from neutral.The attribute shall be greater than lowStep.',
    low_step STRING COMMENT 'Lowest possible tap step position retard from neutral.',
    ltc_flag STRING COMMENT 'Specifies whether or not a TapChanger has load tap changing capabilities.',
    neutral_step STRING COMMENT 'The neutral tap step position for this winding.The attribute shall be equal to or greater than lowStep and equal or less than highStep.It is the step position where the voltage is neutralU when the other terminals of the transformer are at the ratedU.  If there are other tap changers on the transformer those taps are kept constant at their neutralStep.',
    normal_step STRING COMMENT 'The tap step position used in normal network operation for this winding. For a Fixed tap changer indicates the current physical tap setting.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
    step STRING COMMENT 'Tap changer position.Starting step for a steady state solution. Non integer values are allowed to support continuous tap variables. The reasons for continuous value are to support study cases where no discrete tap changer has yet been designed a solution where a narrow voltage band forces the tap step to oscillate or to accommodate for a continuous solution as input.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
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
COMMENT 'Snapshot table for PhaseTapChangerNonLinear Actor state persistence (includes inherited properties from: PhaseTapChanger)'


CREATE TABLE IF NOT EXISTS ontology.grid.phase_tap_changer_symmetrical_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
    high_step STRING COMMENT 'Highest possible tap step position advance from neutral.The attribute shall be greater than lowStep.',
    low_step STRING COMMENT 'Lowest possible tap step position retard from neutral.',
    ltc_flag STRING COMMENT 'Specifies whether or not a TapChanger has load tap changing capabilities.',
    neutral_step STRING COMMENT 'The neutral tap step position for this winding.The attribute shall be equal to or greater than lowStep and equal or less than highStep.It is the step position where the voltage is neutralU when the other terminals of the transformer are at the ratedU.  If there are other tap changers on the transformer those taps are kept constant at their neutralStep.',
    normal_step STRING COMMENT 'The tap step position used in normal network operation for this winding. For a Fixed tap changer indicates the current physical tap setting.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
    step STRING COMMENT 'Tap changer position.Starting step for a steady state solution. Non integer values are allowed to support continuous tap variables. The reasons for continuous value are to support study cases where no discrete tap changer has yet been designed a solution where a narrow voltage band forces the tap step to oscillate or to accommodate for a continuous solution as input.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
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
COMMENT 'Snapshot table for PhaseTapChangerSymmetrical Actor state persistence (includes inherited properties from: PhaseTapChangerNonLinear)'


CREATE TABLE IF NOT EXISTS ontology.grid.phase_tap_changer_tabular_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
    high_step STRING COMMENT 'Highest possible tap step position advance from neutral.The attribute shall be greater than lowStep.',
    low_step STRING COMMENT 'Lowest possible tap step position retard from neutral.',
    ltc_flag STRING COMMENT 'Specifies whether or not a TapChanger has load tap changing capabilities.',
    neutral_step STRING COMMENT 'The neutral tap step position for this winding.The attribute shall be equal to or greater than lowStep and equal or less than highStep.It is the step position where the voltage is neutralU when the other terminals of the transformer are at the ratedU.  If there are other tap changers on the transformer those taps are kept constant at their neutralStep.',
    normal_step STRING COMMENT 'The tap step position used in normal network operation for this winding. For a Fixed tap changer indicates the current physical tap setting.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
    step STRING COMMENT 'Tap changer position.Starting step for a steady state solution. Non integer values are allowed to support continuous tap variables. The reasons for continuous value are to support study cases where no discrete tap changer has yet been designed a solution where a narrow voltage band forces the tap step to oscillate or to accommodate for a continuous solution as input.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
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
COMMENT 'Snapshot table for PhaseTapChangerTabular Actor state persistence (includes inherited properties from: PhaseTapChanger)'


CREATE TABLE IF NOT EXISTS ontology.grid.photo_voltaic_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for PhotoVoltaicUnit Actor state persistence (includes inherited properties from: PowerElectronicsUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.plant_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Plant Actor state persistence (includes inherited properties from: EquipmentContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.post_line_sensor_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for PostLineSensor Actor state persistence (includes inherited properties from: Sensor)'


CREATE TABLE IF NOT EXISTS ontology.grid.potential_transformer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    accuracy_class STRING COMMENT 'PT accuracy classification.',
    nominal_ratio STRING COMMENT 'Nominal ratio between the primary and secondary voltage.',
    pt_class STRING COMMENT 'Potential transformer PT classification covering burden.',
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
COMMENT 'Snapshot table for PotentialTransformer Actor state persistence (includes inherited properties from: Sensor)'


CREATE TABLE IF NOT EXISTS ontology.grid.power_cut_zone_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for PowerCutZone Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.power_electronics_connection_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for PowerElectronicsConnection Actor state persistence (includes inherited properties from: RegulatingCondEq)'


CREATE TABLE IF NOT EXISTS ontology.grid.power_electronics_connection_phase_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for PowerElectronicsConnectionPhase Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.power_electronics_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for PowerElectronicsUnit Actor state persistence (includes inherited properties from: Equipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.power_electronics_wind_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for PowerElectronicsWindUnit Actor state persistence (includes inherited properties from: PowerElectronicsUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.power_system_resource_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for PowerSystemResource Actor state persistence (includes inherited properties from: IdentifiedObject)'


CREATE TABLE IF NOT EXISTS ontology.grid.power_transformer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    is_part_of_generator_unit STRING COMMENT 'Indicates whether the machine is part of a power station unit. Used for short circuit data exchange according to IEC 60909.  It has an impact on how the correction factors are calculated for transformers since the transformer is not necessarily part of a synchronous machine and generating unit. It is not always possible to derive this information from the model. This is why the attribute is necessary.',
    operational_values_considered STRING COMMENT 'It is used to define if the data other attributes related to short circuit data exchange defines long term operational conditions or not. Used for short circuit data exchange according to IEC 60909.',
    vector_group STRING COMMENT 'Vector group of the transformer for protective relaying e.g. Dyn1. For unbalanced transformers this may not be simply determined from the constituent winding connections and phase angle displacements.The vectorGroup string consists of the following components in the order listed high voltage winding connection mid voltage winding connection for three winding transformers phase displacement clock number from 0 to 11  low voltage winding connection phase displacement clock number from 0 to 11.   The winding connections are D delta Y wye YN wye with neutral Z zigzag ZN zigzag with neutral A auto transformer. Upper case means the high voltage lower case mid or low. The high voltage winding always has clock position 0 and is not included in the vector group string.  Some examples YNy0 two winding wye to wye with no phase displacement YNd11 two winding wye to delta with 330 degrees phase displacement YNyn0d5 three winding transformer wye with neutral high voltage wye with neutral mid voltage and no phase displacement delta low voltage with 150 degrees displacement.Phase displacement is defined as the angular difference between the phasors representing the voltages between the neutral point real or imaginary and the corresponding terminals of two windings a positive sequence voltage system being applied to the highvoltage terminals following each other in alphabetical sequence if they are lettered or in numerical sequence if they are numbered the phasors are assumed to rotate in a counterclockwise sense.',
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
COMMENT 'Snapshot table for PowerTransformer Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.prime_mover_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    prime_mover_rating STRING COMMENT 'Rating of prime mover.',
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
COMMENT 'Snapshot table for PrimeMover Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.protected_switch_snapshots (
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
COMMENT 'Snapshot table for ProtectedSwitch Actor state persistence (includes inherited properties from: Switch)'


CREATE TABLE IF NOT EXISTS ontology.grid.protection_equipment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    high_limit STRING COMMENT 'The maximum allowable value.',
    low_limit STRING COMMENT 'The minimum allowable value.',
    power_direction_flag STRING COMMENT 'Direction same as positive active power flow value.',
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
COMMENT 'Snapshot table for ProtectionEquipment Actor state persistence (includes inherited properties from: Equipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.ratio_tap_changer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
    high_step STRING COMMENT 'Highest possible tap step position advance from neutral.The attribute shall be greater than lowStep.',
    low_step STRING COMMENT 'Lowest possible tap step position retard from neutral.',
    ltc_flag STRING COMMENT 'Specifies whether or not a TapChanger has load tap changing capabilities.',
    neutral_step STRING COMMENT 'The neutral tap step position for this winding.The attribute shall be equal to or greater than lowStep and equal or less than highStep.It is the step position where the voltage is neutralU when the other terminals of the transformer are at the ratedU.  If there are other tap changers on the transformer those taps are kept constant at their neutralStep.',
    normal_step STRING COMMENT 'The tap step position used in normal network operation for this winding. For a Fixed tap changer indicates the current physical tap setting.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
    step STRING COMMENT 'Tap changer position.Starting step for a steady state solution. Non integer values are allowed to support continuous tap variables. The reasons for continuous value are to support study cases where no discrete tap changer has yet been designed a solution where a narrow voltage band forces the tap step to oscillate or to accommodate for a continuous solution as input.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
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
COMMENT 'Snapshot table for RatioTapChanger Actor state persistence (includes inherited properties from: TapChanger)'


CREATE TABLE IF NOT EXISTS ontology.grid.recloser_snapshots (
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
COMMENT 'Snapshot table for Recloser Actor state persistence (includes inherited properties from: ProtectedSwitch)'


CREATE TABLE IF NOT EXISTS ontology.grid.registered_distributed_resource_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    distributed_resource_type STRING COMMENT 'The type of resource. Examples include fuel cell flywheel photovoltaic microturbine CHP combined heat power V2G vehicle to grid DES distributed energy storage and others.',
    commercial_op_date STRING COMMENT 'Resource Commercial Operation Date. ',
    dispatchable STRING COMMENT 'Dispatchable indicates whether the resource is dispatchable. This implies that the resource intends to submit Energy bidsoffers or Ancillary Services bidsoffers or selfprovided schedules.',
    last_modified STRING COMMENT 'Indication of the last time this item was modifiedversioned.',
    max_base_self_sched_qty STRING COMMENT 'Maximum base self schedule quantity.',
    max_on_time STRING COMMENT 'Maximum on time after start up.',
    min_off_time STRING COMMENT 'Minimum off time after shut down.',
    min_on_time STRING COMMENT 'Minimum on time after start up.',
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
COMMENT 'Snapshot table for RegisteredDistributedResource Actor state persistence (includes inherited properties from: RegisteredResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.registered_generator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    cold_start_time STRING COMMENT 'Cold start time.',
    combined_cycle_operating_mode STRING COMMENT 'Combined Cycle operating mode.',
    commericial_operation_date STRING COMMENT '',
    hot_int_time STRING COMMENT 'Hottointermediate time Seasonal',
    hot_start_time STRING COMMENT 'Hot start time.',
    int_cold_time STRING COMMENT 'Intermediatetocold time Seasonal',
    int_start_time STRING COMMENT 'Intermediate start time.',
    max_shutdown_time STRING COMMENT 'Maximum time this device can be shut down.',
    max_start_ups_per_day STRING COMMENT 'maximum start ups per day',
    max_weekly_energy STRING COMMENT 'Maximum weekly Energy Seasonal',
    max_weekly_starts STRING COMMENT 'Maximum weekly starts seasonal parameter',
    pump_min_down_time STRING COMMENT 'The minimum down time for the pump in a pump storage unit.',
    pump_min_up_time STRING COMMENT 'The minimum up time aspect for the pump in a pump storage unit',
    pump_shutdown_cost STRING COMMENT 'The cost to shutdown a pump during the pump aspect of a pump storage unit.',
    pump_shutdown_time STRING COMMENT 'The shutdown time minutes of the pump aspect of a pump storage unit.',
    pumping_factor STRING COMMENT 'Pumping factor for pump storage units conversion factor between generating and pumping.',
    resource_sub_type STRING COMMENT 'Unit sub type used by Settlements or scheduling application. Application use of the unit sub type may define the necessary types as applicable.',
    river_system STRING COMMENT 'River System the Resource is tied to.',
    commercial_op_date STRING COMMENT 'Resource Commercial Operation Date. ',
    dispatchable STRING COMMENT 'Dispatchable indicates whether the resource is dispatchable. This implies that the resource intends to submit Energy bidsoffers or Ancillary Services bidsoffers or selfprovided schedules.',
    last_modified STRING COMMENT 'Indication of the last time this item was modifiedversioned.',
    max_base_self_sched_qty STRING COMMENT 'Maximum base self schedule quantity.',
    max_on_time STRING COMMENT 'Maximum on time after start up.',
    min_off_time STRING COMMENT 'Minimum off time after shut down.',
    min_on_time STRING COMMENT 'Minimum on time after start up.',
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
COMMENT 'Snapshot table for RegisteredGenerator Actor state persistence (includes inherited properties from: RegisteredResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.registered_inter_tie_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    min_hourly_block_limit STRING COMMENT 'The registered upper bound of minimum hourly block for an InterTie Resource.',
    commercial_op_date STRING COMMENT 'Resource Commercial Operation Date. ',
    dispatchable STRING COMMENT 'Dispatchable indicates whether the resource is dispatchable. This implies that the resource intends to submit Energy bidsoffers or Ancillary Services bidsoffers or selfprovided schedules.',
    last_modified STRING COMMENT 'Indication of the last time this item was modifiedversioned.',
    max_base_self_sched_qty STRING COMMENT 'Maximum base self schedule quantity.',
    max_on_time STRING COMMENT 'Maximum on time after start up.',
    min_off_time STRING COMMENT 'Minimum off time after shut down.',
    min_on_time STRING COMMENT 'Minimum on time after start up.',
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
COMMENT 'Snapshot table for RegisteredInterTie Actor state persistence (includes inherited properties from: RegisteredResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.registered_load_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    block_load_transfer STRING COMMENT 'Emergency operating procedure  Flag to indicate that the Resource is Block Load pseudo resource.',
    dynamically_scheduled_load_resource STRING COMMENT 'Flag to indicate that a Load Resource is part of a DSR Load',
    dynamically_scheduled_qualification STRING COMMENT 'Qualification status used for DSR qualification.',
    commercial_op_date STRING COMMENT 'Resource Commercial Operation Date. ',
    dispatchable STRING COMMENT 'Dispatchable indicates whether the resource is dispatchable. This implies that the resource intends to submit Energy bidsoffers or Ancillary Services bidsoffers or selfprovided schedules.',
    last_modified STRING COMMENT 'Indication of the last time this item was modifiedversioned.',
    max_base_self_sched_qty STRING COMMENT 'Maximum base self schedule quantity.',
    max_on_time STRING COMMENT 'Maximum on time after start up.',
    min_off_time STRING COMMENT 'Minimum off time after shut down.',
    min_on_time STRING COMMENT 'Minimum on time after start up.',
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
COMMENT 'Snapshot table for RegisteredLoad Actor state persistence (includes inherited properties from: RegisteredResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.registered_resource_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    commercial_op_date STRING COMMENT 'Resource Commercial Operation Date. ',
    dispatchable STRING COMMENT 'Dispatchable indicates whether the resource is dispatchable. This implies that the resource intends to submit Energy bidsoffers or Ancillary Services bidsoffers or selfprovided schedules.',
    last_modified STRING COMMENT 'Indication of the last time this item was modifiedversioned.',
    max_base_self_sched_qty STRING COMMENT 'Maximum base self schedule quantity.',
    max_on_time STRING COMMENT 'Maximum on time after start up.',
    min_off_time STRING COMMENT 'Minimum off time after shut down.',
    min_on_time STRING COMMENT 'Minimum on time after start up.',
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
COMMENT 'Snapshot table for RegisteredResource Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.regulating_cond_eq_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for RegulatingCondEq Actor state persistence (includes inherited properties from: EnergyConnection)'


CREATE TABLE IF NOT EXISTS ontology.grid.regulating_control_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    discrete STRING COMMENT 'The regulation is performed in a discrete mode. This applies to equipment with discrete controls e.g. tap changers and shunt compensators.',
    enabled STRING COMMENT 'The flag tells if regulation is enabled.',
    max_allowed_target_value STRING COMMENT 'Maximum allowed target value RegulatingControl.targetValue.',
    min_allowed_target_value STRING COMMENT 'Minimum allowed target value RegulatingControl.targetValue.',
    target_deadband STRING COMMENT 'This is a deadband used with discrete control to avoid excessive update of controls like tap changers and shunt compensator banks while regulating.  The units of those appropriate for the mode. The attribute shall be a positive value or zero. If RegulatingControl.discrete is set to false the RegulatingControl.targetDeadband is to be ignored.Note that for instance if the targetValue is 100 kV and the targetDeadband is 2 kV the range is from 99 to 101 kV.',
    target_value STRING COMMENT 'The target value specified for case input.   This value can be used for the target value without the use of schedules. The value has the units appropriate to the mode attribute.',
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
COMMENT 'Snapshot table for RegulatingControl Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.remedial_action_scheme_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    armed STRING COMMENT 'The status of the class set by operation or by signal. Optional field that will override other status fields.',
    normal_armed STRING COMMENT 'The defaultnormal value used when other active signalvalues are missing.',
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
COMMENT 'Snapshot table for RemedialActionScheme Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.remote_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for RemoteUnit Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.reservoir_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    energy_storage_rating STRING COMMENT 'The reservoirs energy storage rating in energy for given head conditions.',
    river_outlet_works STRING COMMENT 'River outlet works for riparian right releases or other purposes.',
    spill_way_gate_type STRING COMMENT 'Type of spillway gate including parameters.',
    spillway_capacity STRING COMMENT 'The flow capacity of the spillway in cubic meters per second.',
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
COMMENT 'Snapshot table for Reservoir Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.rotating_machine_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    rated_power_factor STRING COMMENT 'Power factor nameplate data. It is primarily used for short circuit data exchange according to IEC 60909. The attribute cannot be a negative value.',
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for RotatingMachine Actor state persistence (includes inherited properties from: RegulatingCondEq)'


CREATE TABLE IF NOT EXISTS ontology.grid.svc_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    grounded STRING COMMENT 'Used for Yn and Zn connections. True if the neutral is solidly grounded.',
    maximum_sections STRING COMMENT 'The maximum number of sections that may be switched in. ',
    normal_sections STRING COMMENT 'The normal number of sections switched in. The value shall be between zero and ShuntCompensator.maximumSections.',
    sections STRING COMMENT 'Shunt compensator sections in use. Starting value for steady state solution. The attribute shall be a positive value or zero. Non integer values are allowed to support continuous variables. The reasons for continuous value are to support study cases where no discrete shunt compensators has yet been designed a solutions where a narrow voltage band force the sections to oscillate or accommodate for a continuous solution as input. For LinearShuntConpensator the value shall be between zero and ShuntCompensator.maximumSections. At value zero the shunt compensator conductance and admittance is zero. Linear interpolation of conductance and admittance between the previous and next integer section is applied in case of noninteger values.For NonlinearShuntCompensators shall only be set to one of the NonlinearShuntCompenstorPoint.sectionNumber. There is no interpolation between NonlinearShuntCompenstorPoints.',
    switch_on_count STRING COMMENT 'The switch on count since the capacitor count was last reset or initialized.',
    switch_on_date STRING COMMENT 'The date and time when the capacitor bank was last switched on.',
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for SVC Actor state persistence (includes inherited properties from: ShuntCompensator)'


CREATE TABLE IF NOT EXISTS ontology.grid.sectionaliser_snapshots (
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
COMMENT 'Snapshot table for Sectionaliser Actor state persistence (includes inherited properties from: Switch)'


CREATE TABLE IF NOT EXISTS ontology.grid.sensor_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Sensor Actor state persistence (includes inherited properties from: AuxiliaryEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.series_compensator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    varistor_present STRING COMMENT 'Describe if a metal oxide varistor mov for over voltage protection is configured in parallel with the series compensator. It is used for short circuit calculations.',
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
COMMENT 'Snapshot table for SeriesCompensator Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.shunt_compensator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    grounded STRING COMMENT 'Used for Yn and Zn connections. True if the neutral is solidly grounded.',
    maximum_sections STRING COMMENT 'The maximum number of sections that may be switched in. ',
    normal_sections STRING COMMENT 'The normal number of sections switched in. The value shall be between zero and ShuntCompensator.maximumSections.',
    sections STRING COMMENT 'Shunt compensator sections in use. Starting value for steady state solution. The attribute shall be a positive value or zero. Non integer values are allowed to support continuous variables. The reasons for continuous value are to support study cases where no discrete shunt compensators has yet been designed a solutions where a narrow voltage band force the sections to oscillate or accommodate for a continuous solution as input. For LinearShuntConpensator the value shall be between zero and ShuntCompensator.maximumSections. At value zero the shunt compensator conductance and admittance is zero. Linear interpolation of conductance and admittance between the previous and next integer section is applied in case of noninteger values.For NonlinearShuntCompensators shall only be set to one of the NonlinearShuntCompenstorPoint.sectionNumber. There is no interpolation between NonlinearShuntCompenstorPoints.',
    switch_on_count STRING COMMENT 'The switch on count since the capacitor count was last reset or initialized.',
    switch_on_date STRING COMMENT 'The date and time when the capacitor bank was last switched on.',
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for ShuntCompensator Actor state persistence (includes inherited properties from: RegulatingCondEq)'


CREATE TABLE IF NOT EXISTS ontology.grid.shunt_compensator_control_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    branch_direct STRING COMMENT 'For VAR amp or power factor locally controlled shunt impedances the flow direction in out.',
    local_off_level STRING COMMENT 'Upper control setting.',
    local_on_level STRING COMMENT 'Lower control setting.',
    local_override STRING COMMENT 'True if the locally controlled capacitor has voltage override capability.',
    max_switch_operation_count STRING COMMENT 'IdmsShuntImpedanceData.maxNumSwitchOps.',
    normal_open STRING COMMENT 'True if open is normal status for a fixed capacitor bank otherwise normal status is closed.',
    reg_branch STRING COMMENT 'For VAR amp or power factor locally controlled shunt impedances the index of the regulation branch.',
    reg_branch_end STRING COMMENT 'For VAR amp or power factor locally controlled shunt impedances the end of the branch that is regulated. The field has the following values from side to side and tertiary only if the branch is a transformer.',
    v_reg_line_line STRING COMMENT 'True if regulated voltages are measured line to line otherwise they are measured line to ground.',
    discrete STRING COMMENT 'The regulation is performed in a discrete mode. This applies to equipment with discrete controls e.g. tap changers and shunt compensators.',
    enabled STRING COMMENT 'The flag tells if regulation is enabled.',
    max_allowed_target_value STRING COMMENT 'Maximum allowed target value RegulatingControl.targetValue.',
    min_allowed_target_value STRING COMMENT 'Minimum allowed target value RegulatingControl.targetValue.',
    target_deadband STRING COMMENT 'This is a deadband used with discrete control to avoid excessive update of controls like tap changers and shunt compensator banks while regulating.  The units of those appropriate for the mode. The attribute shall be a positive value or zero. If RegulatingControl.discrete is set to false the RegulatingControl.targetDeadband is to be ignored.Note that for instance if the targetValue is 100 kV and the targetDeadband is 2 kV the range is from 99 to 101 kV.',
    target_value STRING COMMENT 'The target value specified for case input.   This value can be used for the target value without the use of schedules. The value has the units appropriate to the mode attribute.',
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
COMMENT 'Snapshot table for ShuntCompensatorControl Actor state persistence (includes inherited properties from: RegulatingControl)'


CREATE TABLE IF NOT EXISTS ontology.grid.shunt_compensator_phase_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    maximum_sections STRING COMMENT 'The maximum number of sections that may be switched in for this phase. ',
    normal_sections STRING COMMENT 'For the capacitor phase the normal number of sections switched in. The value shall be between zero and ShuntCompensatorPhase.maximumSections.',
    sections STRING COMMENT 'Shunt compensator sections in use. Starting value for steady state solution. The attribute shall be a positive value or zero. Non integer values are allowed to support continuous variables. The reasons for continuous value are to support study cases where no discrete shunt compensators has yet been designed a solutions where a narrow voltage band force the sections to oscillate or accommodate for a continuous solution as input.For LinearShuntConpensator the value shall be between zero and ShuntCompensatorPhase.maximumSections. At value zero the shunt compensator conductance and admittance is zero. Linear interpolation of conductance and admittance between the previous and next integer section is applied in case of noninteger values.For NonlinearShuntCompensators shall only be set to one of the NonlinearShuntCompenstorPhasePoint.sectionNumber. There is no interpolation between NonlinearShuntCompenstorPhasePoints.',
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
COMMENT 'Snapshot table for ShuntCompensatorPhase Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.solar_generating_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    long_pf STRING COMMENT 'Generating unit long term economic participation factor.',
    normal_pf STRING COMMENT 'Generating unit economic participation factor.  The sum of the participation factors across generating units does not have to sum to one.  It is used for representing distributed slack participation factor. The attribute shall be a positive value or zero.',
    penalty_factor STRING COMMENT 'Defined as 1   1  Incremental Transmission Loss with the Incremental Transmission Loss expressed as a plus or minus value. The typical range of penalty factors is 0.9 to 1.1.',
    short_pf STRING COMMENT 'Generating unit short term economic participation factor.',
    tie_line_pf STRING COMMENT 'Generating unit economic participation factor.',
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
COMMENT 'Snapshot table for SolarGeneratingUnit Actor state persistence (includes inherited properties from: GeneratingUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.static_var_compensator_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for StaticVarCompensator Actor state persistence (includes inherited properties from: RegulatingCondEq)'


CREATE TABLE IF NOT EXISTS ontology.grid.station_supply_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    customer_count STRING COMMENT 'Number of individual customers represented by this demand.',
    grounded STRING COMMENT 'Used for Yn and Zn connections. True if the neutral is solidly grounded.',
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
COMMENT 'Snapshot table for StationSupply Actor state persistence (includes inherited properties from: EnergyConsumer)'


CREATE TABLE IF NOT EXISTS ontology.grid.steam_supply_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    steam_supply_rating STRING COMMENT 'Rating of steam supply.',
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
COMMENT 'Snapshot table for SteamSupply Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.steam_turbine_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    shaft1_power_hp STRING COMMENT 'Fraction of power from shaft 1 high pressure turbine output.',
    shaft1_power_ip STRING COMMENT 'Fraction of power from shaft 1 intermediate pressure turbine output.',
    shaft1_power_lp1 STRING COMMENT 'Fraction of power from shaft 1 first low pressure turbine output.',
    shaft1_power_lp2 STRING COMMENT 'Fraction of power from shaft 1 second low pressure turbine output.',
    shaft2_power_hp STRING COMMENT 'Fraction of power from shaft 2 high pressure turbine output.',
    shaft2_power_ip STRING COMMENT 'Fraction of power from shaft 2 intermediate pressure turbine output.',
    shaft2_power_lp1 STRING COMMENT 'Fraction of power from shaft 2 first low pressure turbine output.',
    shaft2_power_lp2 STRING COMMENT 'Fraction of power from shaft 2 second low pressure turbine output.',
    prime_mover_rating STRING COMMENT 'Rating of prime mover.',
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
COMMENT 'Snapshot table for SteamTurbine Actor state persistence (includes inherited properties from: PrimeMover)'


CREATE TABLE IF NOT EXISTS ontology.grid.sub_control_area_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    area_short_name STRING COMMENT 'Market area short name which is the regulation zone. It references AGC regulation zone name.',
    constant_coefficient STRING COMMENT 'Loss estimate constant coefficient',
    linear_coefficient STRING COMMENT 'Loss estimate linear coefficient',
    max_self_sched_mw STRING COMMENT 'Maximum amount of self schedule MWs allowed for an embedded control area.',
    min_self_sched_mw STRING COMMENT 'Minimum amount of self schedule MW allowed for an embedded control area.',
    quadratic_coefficient STRING COMMENT 'Loss estimate quadratic coefficient',
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
COMMENT 'Snapshot table for SubControlArea Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.subcritical_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_error_bias_p STRING COMMENT 'Active power error bias ratio.',
    control_ic STRING COMMENT 'Integral constant.',
    control_pc STRING COMMENT 'Proportional constant.',
    control_peb STRING COMMENT 'Pressure error bias ratio.',
    control_tc STRING COMMENT 'Time constant.',
    feed_water_ig STRING COMMENT 'Feedwater integral gain ratio.',
    feed_water_pg STRING COMMENT 'Feedwater proportional gain ratio.',
    max_error_rate_p STRING COMMENT 'Active power maximum error rate limit.',
    min_error_rate_p STRING COMMENT 'Active power minimum error rate limit.',
    pressure_ctrl_dg STRING COMMENT 'Pressure control derivative gain ratio.',
    pressure_ctrl_ig STRING COMMENT 'Pressure control integral gain ratio.',
    pressure_ctrl_pg STRING COMMENT 'Pressure control proportional gain ratio.',
    pressure_feedback STRING COMMENT 'Pressure feedback indicator.',
    super_heater1_capacity STRING COMMENT 'Drumprimary superheater capacity.',
    super_heater2_capacity STRING COMMENT 'Secondary superheater capacity.',
    super_heater_pipe_pd STRING COMMENT 'Superheater pipe pressure drop constant.',
    steam_supply_rating STRING COMMENT 'Rating of steam supply.',
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
COMMENT 'Snapshot table for Subcritical Actor state persistence (includes inherited properties from: FossilSteamSupply)'


CREATE TABLE IF NOT EXISTS ontology.grid.substation_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for Substation Actor state persistence (includes inherited properties from: EquipmentContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.supercritical_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_error_bias_p STRING COMMENT 'Active power error bias ratio.',
    control_ic STRING COMMENT 'Integral constant.',
    control_pc STRING COMMENT 'Proportional constant.',
    control_peb STRING COMMENT 'Pressure error bias ratio.',
    control_tc STRING COMMENT 'Time constant.',
    feed_water_ig STRING COMMENT 'Feedwater integral gain ratio.',
    feed_water_pg STRING COMMENT 'Feedwater proportional gain ratio.',
    max_error_rate_p STRING COMMENT 'Active power maximum error rate limit.',
    min_error_rate_p STRING COMMENT 'Active power minimum error rate limit.',
    pressure_ctrl_dg STRING COMMENT 'Pressure control derivative gain ratio.',
    pressure_ctrl_ig STRING COMMENT 'Pressure control integral gain ratio.',
    pressure_ctrl_pg STRING COMMENT 'Pressure control proportional gain ratio.',
    pressure_feedback STRING COMMENT 'Pressure feedback indicator.',
    super_heater1_capacity STRING COMMENT 'Drumprimary superheater capacity.',
    super_heater2_capacity STRING COMMENT 'Secondary superheater capacity.',
    super_heater_pipe_pd STRING COMMENT 'Superheater pipe pressure drop constant.',
    steam_supply_rating STRING COMMENT 'Rating of steam supply.',
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
COMMENT 'Snapshot table for Supercritical Actor state persistence (includes inherited properties from: FossilSteamSupply)'


CREATE TABLE IF NOT EXISTS ontology.grid.surge_arrester_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for SurgeArrester Actor state persistence (includes inherited properties from: AuxiliaryEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.switch_snapshots (
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
COMMENT 'Snapshot table for Switch Actor state persistence (includes inherited properties from: ConductingEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.switch_phase_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    closed STRING COMMENT 'The attribute tells if the switch is considered closed when used as input to topology processing.',
    normal_open STRING COMMENT 'Used in cases when no Measurement for the status value is present. If the SwitchPhase has a status measurement the Discrete.normalValue is expected to match with this value.',
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
COMMENT 'Snapshot table for SwitchPhase Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.synchrocheck_relay_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    high_limit STRING COMMENT 'The maximum allowable value.',
    low_limit STRING COMMENT 'The minimum allowable value.',
    power_direction_flag STRING COMMENT 'Direction same as positive active power flow value.',
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
COMMENT 'Snapshot table for SynchrocheckRelay Actor state persistence (includes inherited properties from: ProtectionEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.synchronous_machine_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    coolant_condition STRING COMMENT 'Temperature or pressure of coolant medium.',
    earthing STRING COMMENT 'Indicates whether or not the generator is earthed. Used for short circuit data exchange according to IEC 60909.',
    mu STRING COMMENT 'Factor to calculate the breaking current Section 4.5.2.1 in IEC 609090.Used only for single fed short circuit on a generator Section 4.3.4.2. in IEC 609090.',
    reference_priority STRING COMMENT 'Priority of unit for use as powerflow voltage phase angle reference bus selection. 0  don t care default 1  highest priority. 2 is less than 1 and so on.',
    rated_power_factor STRING COMMENT 'Power factor nameplate data. It is primarily used for short circuit data exchange according to IEC 60909. The attribute cannot be a negative value.',
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
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
COMMENT 'Snapshot table for SynchronousMachine Actor state persistence (includes inherited properties from: RotatingMachine)'


CREATE TABLE IF NOT EXISTS ontology.grid.tcp_access_point_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    keep_alive_time STRING COMMENT 'Indicates the default interval at which TCP will check if the TCP connection is still valid.',
    port STRING COMMENT 'This value is only needed to be specified for called nodes e.g. those that respond to a TCP.Open request.This value specifies the TCP port to be used. Well known and registered ports are preferred and can be found at httpwww.iana.orgassignmentsservicenamesportnumbersservicenamesportnumbers.xhtmlFor IEC 608706 TASE.2 e.g. ICCP and IEC 61850 the value used shall be 102 for nonTLS protected exchanges. The value shall be 3782 for TLS transported ICCP and 61850 exchanges.',
    address STRING COMMENT 'Is the dotted decimal IP Address resolve the IP address. The format is controlled by the value of the addressType.',
    gateway STRING COMMENT 'Is the dotted decimal IPAddress of the first hop router.  Format is controlled by the addressType.',
    subnet STRING COMMENT 'This is the IP subnet mask which controls the local vs nonlocal routing.',
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
COMMENT 'Snapshot table for TCPAccessPoint Actor state persistence (includes inherited properties from: IPAccessPoint)'


CREATE TABLE IF NOT EXISTS ontology.grid.tap_changer_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    control_enabled STRING COMMENT 'Specifies the regulation status of the equipment.  True is regulating false is not regulating.',
    high_step STRING COMMENT 'Highest possible tap step position advance from neutral.The attribute shall be greater than lowStep.',
    low_step STRING COMMENT 'Lowest possible tap step position retard from neutral.',
    ltc_flag STRING COMMENT 'Specifies whether or not a TapChanger has load tap changing capabilities.',
    neutral_step STRING COMMENT 'The neutral tap step position for this winding.The attribute shall be equal to or greater than lowStep and equal or less than highStep.It is the step position where the voltage is neutralU when the other terminals of the transformer are at the ratedU.  If there are other tap changers on the transformer those taps are kept constant at their neutralStep.',
    normal_step STRING COMMENT 'The tap step position used in normal network operation for this winding. For a Fixed tap changer indicates the current physical tap setting.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
    step STRING COMMENT 'Tap changer position.Starting step for a steady state solution. Non integer values are allowed to support continuous tap variables. The reasons for continuous value are to support study cases where no discrete tap changer has yet been designed a solution where a narrow voltage band forces the tap step to oscillate or to accommodate for a continuous solution as input.The attribute shall be equal to or greater than lowStep and equal to or less than highStep.',
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
COMMENT 'Snapshot table for TapChanger Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.tap_changer_control_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    line_drop_compensation STRING COMMENT 'If true the line drop compensation is to be applied.',
    discrete STRING COMMENT 'The regulation is performed in a discrete mode. This applies to equipment with discrete controls e.g. tap changers and shunt compensators.',
    enabled STRING COMMENT 'The flag tells if regulation is enabled.',
    max_allowed_target_value STRING COMMENT 'Maximum allowed target value RegulatingControl.targetValue.',
    min_allowed_target_value STRING COMMENT 'Minimum allowed target value RegulatingControl.targetValue.',
    target_deadband STRING COMMENT 'This is a deadband used with discrete control to avoid excessive update of controls like tap changers and shunt compensator banks while regulating.  The units of those appropriate for the mode. The attribute shall be a positive value or zero. If RegulatingControl.discrete is set to false the RegulatingControl.targetDeadband is to be ignored.Note that for instance if the targetValue is 100 kV and the targetDeadband is 2 kV the range is from 99 to 101 kV.',
    target_value STRING COMMENT 'The target value specified for case input.   This value can be used for the target value without the use of schedules. The value has the units appropriate to the mode attribute.',
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
COMMENT 'Snapshot table for TapChangerControl Actor state persistence (includes inherited properties from: RegulatingControl)'


CREATE TABLE IF NOT EXISTS ontology.grid.thermal_generating_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    long_pf STRING COMMENT 'Generating unit long term economic participation factor.',
    normal_pf STRING COMMENT 'Generating unit economic participation factor.  The sum of the participation factors across generating units does not have to sum to one.  It is used for representing distributed slack participation factor. The attribute shall be a positive value or zero.',
    penalty_factor STRING COMMENT 'Defined as 1   1  Incremental Transmission Loss with the Incremental Transmission Loss expressed as a plus or minus value. The typical range of penalty factors is 0.9 to 1.1.',
    short_pf STRING COMMENT 'Generating unit short term economic participation factor.',
    tie_line_pf STRING COMMENT 'Generating unit economic participation factor.',
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
COMMENT 'Snapshot table for ThermalGeneratingUnit Actor state persistence (includes inherited properties from: GeneratingUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.transformer_tank_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for TransformerTank Actor state persistence (includes inherited properties from: Equipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.transmission_corridor_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for TransmissionCorridor Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.transmission_right_of_way_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for TransmissionRightOfWay Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.voltage_control_zone_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for VoltageControlZone Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.voltage_level_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for VoltageLevel Actor state persistence (includes inherited properties from: EquipmentContainer)'


CREATE TABLE IF NOT EXISTS ontology.grid.vs_converter_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    max_modulation_index STRING COMMENT 'The maximum quotient between the AC converter voltage Uc and DC voltage Ud. A factor typically less than 1. It is converters configuration data used in power flow.',
    target_pw_mfactor STRING COMMENT 'Magnitude of pulsemodulation factor. The attribute shall be a positive value.',
    target_power_factor_pcc STRING COMMENT 'Power factor target at the AC side at point of common coupling. The attribute shall be a positive value.',
    number_of_valves STRING COMMENT 'Number of valves in the converter. Used in loss calculations.',
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
COMMENT 'Snapshot table for VsConverter Actor state persistence (includes inherited properties from: ACDCConverter)'


CREATE TABLE IF NOT EXISTS ontology.grid.wave_trap_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for WaveTrap Actor state persistence (includes inherited properties from: AuxiliaryEquipment)'


CREATE TABLE IF NOT EXISTS ontology.grid.weather_station_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for WeatherStation Actor state persistence (includes inherited properties from: PowerSystemResource)'


CREATE TABLE IF NOT EXISTS ontology.grid.wind_generating_unit_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    long_pf STRING COMMENT 'Generating unit long term economic participation factor.',
    normal_pf STRING COMMENT 'Generating unit economic participation factor.  The sum of the participation factors across generating units does not have to sum to one.  It is used for representing distributed slack participation factor. The attribute shall be a positive value or zero.',
    penalty_factor STRING COMMENT 'Defined as 1   1  Incremental Transmission Loss with the Incremental Transmission Loss expressed as a plus or minus value. The typical range of penalty factors is 0.9 to 1.1.',
    short_pf STRING COMMENT 'Generating unit short term economic participation factor.',
    tie_line_pf STRING COMMENT 'Generating unit economic participation factor.',
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
COMMENT 'Snapshot table for WindGeneratingUnit Actor state persistence (includes inherited properties from: GeneratingUnit)'


CREATE TABLE IF NOT EXISTS ontology.grid.wire_segment_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
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
COMMENT 'Snapshot table for WireSegment Actor state persistence (includes inherited properties from: Conductor)'


CREATE TABLE IF NOT EXISTS ontology.grid.wire_segment_phase_snapshots (
actor_id STRING NOT NULL,
owl_class_uri STRING NOT NULL,
sequence BIGINT NOT NULL,
timestamp TIMESTAMP NOT NULL,
snapshot_id STRING,
    sequence_number STRING COMMENT 'Number designation for this wire segment phase. Each wire segment phase within a wire segment should have a unique sequence number.',
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
COMMENT 'Snapshot table for WireSegmentPhase Actor state persistence (includes inherited properties from: PowerSystemResource)'
