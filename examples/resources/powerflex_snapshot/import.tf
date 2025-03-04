/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.
Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://mozilla.org/MPL/2.0/
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//Gather Snapshot existing volumes
data "powerflex_volume" "Snapshot" {
    filter{
        name = ["test-snapshot-import"]
    }
}

//Gather volume linked to existing snapshot
data "powerflex_volume" "Volume" {
    filter{
        name = ["test-vol-import"]
    }
}

//Import all snapshots from volume
import {
    to = powerflex_snapshot.import_test_snapshot
    id = data.powerflex_volume.Snapshot.volumes[0].id
}

//Add them to terraform state
resource "powerflex_snapshot" "import_test_snapshot"{
    name = data.powerflex_volume.Snapshot.volumes[0].name
    volume_id = data.powerflex_volume.Volume.volumes[0].id
}