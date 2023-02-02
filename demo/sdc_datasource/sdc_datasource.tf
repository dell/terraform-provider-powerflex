
# # -----------------------------------------------------------------------------------
# # Read all sdcs if id is blank, otherwise reads all sdcs
# # -----------------------------------------------------------------------------------
    # name is optional if empty then will return all sdc
    # sdcid is optional if empty then will return all sdc
    # sdcid and name both are empty then will return all sdc
data "powerflex_sdc" "selected" {
    # id = "c423b09800000003"
    # name = ""
}

# # Returns all sdcs matching criteria
output "allsdcresult" {
  value = data.powerflex_sdc.selected
}
# # -----------------------------------------------------------------------------------

