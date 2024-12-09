# convert-mpf-inventory

Helping out VPC to take data from MPF, which they bought. Download all iamges and extract data to save it in a format they are able to use.

There been a manual changes around, ex: Extracting collections and their IDs from the Website was a late idea. Becomes some commenting and un-commenting to get everything to work as expected.

## Features

- There is a duplicateCheck.json that is being created, which the tool checks to see if the item has been generated or not, making sure there is no duplications
- Missing SKU, generating a new SKU and it should awalys be unique
- The OutPut JSON is a debugger in a way, to see what happen and able to trace for files and images saved.

## TO DO

- When exporting CSV and JSON, should add **-1** when a file with the same name exist
- Add date to the JSON exports
