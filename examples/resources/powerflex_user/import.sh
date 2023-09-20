# /*
# Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.
# Licensed under the Mozilla Public License Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#     http://mozilla.org/MPL/2.0/
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
# */

# import user by it's id
terraform import powerflex_user.user_import_by_id "<id>"

# import user by it's id - alternative approach by prefixing it with "id:"
terraform import powerflex_user.user_import_by_id "<id:id_of_the_user>"

# import user by it's name
terraform import powerflex_user.user_import_by_name "<name:name_of_the_user>"
