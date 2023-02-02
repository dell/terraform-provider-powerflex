
# # -----------------------------------------------------------------------------------
# # Read all volumes if completely blank, otherwise reads specific volume based on id or name
# # -----------------------------------------------------------------------------------
    # name is optional , if mentioned then will retrieve the specific volume with that name
    # id is optional , if mentioned then will retrieve the specific volume with that id
    # id and name both are empty then will return all sdc

data "powerflex_snapshotpolicy" "sp" {

    #name = "sample_snap_policy_1"
    #id = "15ad99b900000001"
}

output "spResult" {
  value = data.powerflex_snapshotpolicy.sp.snapshotpolicies
}

