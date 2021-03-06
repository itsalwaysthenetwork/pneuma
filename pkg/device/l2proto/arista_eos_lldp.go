package l2proto

const EOS_LLDP_TEMPLATE string = `
Value NEIGHBOR (\S+)
Value LOCAL_INTERFACE (\S+)
Value NEIGHBOR_INTERFACE (\S+)

Start
  ^Port.*TTL -> LLDP

LLDP
  ^${LOCAL_INTERFACE}\s+${NEIGHBOR}\s+${NEIGHBOR_INTERFACE}\s+.* -> Record
`
