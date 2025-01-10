---
page_title: "Datasource Filtering"
title: "Datasource Filtering"
linkTitle: "Datasource Filtering"
---

<!--
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

## Datasource Filtering Overview

*This describes how filtering works in the Powerflex Provider with a datasource filter block*

The input below will be used showcase different scenarios and how they will work in our provider:
```
[
  {
   id = "id-1"
   field = false
   count = 1
  },
  {
   id = "id-2"
   field = true
   count = 2
  },
  {
   id = "id-3"
   field = true
   count = 3
  }
]
```

## 1. If a single filter field is set, filter and only return those values

**Config:**
```
data "powerflex_example_datasource" "exampleFilter" {
  filter {
   id = ["id-1", "id-2"]
  }
}
```
**Output:**
```
[
  {
   id = "id-1"
   field = false
   count = 1
  },
  {
   id = "id-2"
   field = true
   count = 2
  }
]
```

## 2. If multiple filters are set then it is the intersection of those filters:

**Config:**
```
data "powerflex_example_datasource" "exampleFilter" {
  filter {
   id = ["id-1", "id-2"]
   count = [2]
  }
}
```
**Output:**
```
[
  {
   id = "id-2"
   field = true
   count = 2
  }
]
```

## 3. If there are no intersection then the output will be empty:

**Config:**
```
data "powerflex_example_datasource" "exampleFilter" {
  filter {
   id = ["id-1", "id-2"]
   count = [3]
  }
}
```
**Output:**
```
[]
```

## 4. If the filter value is invalid then the output will be empty:
**Config:**
```
data "powerflex_example_datasource" "exampleFilter" {
  filter {
   id = ["invaid-id"]
  }
}
```
**Output:**
```
[]
```

## 5. If the filter value is using regular expressions:
**Config:**
```
data "powerflex_example_datasource" "exampleFilter" {
  filter {
   id = ["^id-[1-2]$"]
  }
}
```
**Output:**
```
[
   {
   id = "id-1"
   field = false
   count = 1
  },
  {
   id = "id-2"
   field = true
   count = 2
  }
]
```