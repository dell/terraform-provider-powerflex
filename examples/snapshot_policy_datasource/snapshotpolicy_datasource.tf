# commands to run this tf file : terraform init && terraform apply --auto-approve
# Reads snapshot policy either by name or by id , if provided
# If both name and id is not provided , then it reads all the snapshot policies
# id and name can't be given together to fetch the snapshot policy

data "powerflex_snapshot_policy" "sp" {

    #name = "sample_snap_policy_1"
    id = "15ad99b900000001"
}

output "spResult" {
  value = data.powerflex_snapshot_policy.sp.snapshotpolicies
}

