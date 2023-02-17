# VsCode Extensions - Osquery

Use Osquery to gather data about the installed VSCode extensions on a host.

## Usage
For testing, you can load the extension with `osqueryi`.

By default, osquery does not want to load extensions not owned by root. You can either change the ownership of vscode_extensions.ext to root, or run osquery with the `--allow_unsafe` flag.

```bash
osqueryi --extension /path/to/vscode_extensions.ext
```

For production deployment, you should refer to the [osquery documentation](https://osquery.readthedocs.io/en/stable/deployment/extensions/).

## Table structure

|COLUMN|TYPE|DESCRIPTION|
|---|---|---|
|name|TEXT|Extension name|
|category|TEXT|The author supplied category of the VSCode extension|
|description|TEXT|The author supplied description of the VSCode extension|
|display_name|TEXT|The user-friendly name of the VSCode extension|
|license|TEXT|License provided by the author|
|path|TEXT|The path where the VSCode extension is installed|
|url|TEXT|The repository url of the VSCode extension|
|version|TEXT|The text representation of the version|
|extension_id|TEXT|Unique extension id|
|identifier|TEXT|Name representation of the identifier|
|publisher|TEXT|Author name|
|publisher_id|TEXT|Author ID|
|installed_at|BIGINT|The date when the extension was installed|
|user|TEXT|User that installed the extension|

## Example

```sql
osquery> select name, category, display_name, version, extension_id, publisher, installed_at from vscode_extensions limit 2;
+-----------+------------------------------------------+---------------------------------+---------+--------------------------------------+---------------------+---------------+
| name      | category                                 | display_name                    | version | extension_id                         | publisher           | installed_at  |
+-----------+------------------------------------------+---------------------------------+---------+--------------------------------------+---------------------+---------------+
| doxdocgen | Other                                    | Doxygen Documentation Generator | 1.4.0   | da7e26d5-d57c-4742-ab47-d77fb189e195 | Christoph Schlosser | 1655372441714 |
| sysmon    | Other - Programming Languages - Snippets | Sysmon                          | 1.8.3   | 0b631589-0c5b-445d-aeab-c1bda6b98946 | DarkOperator        | 1665477449776 |
+-----------+------------------------------------------+---------------------------------+---------+--------------------------------------+---------------------+---------------+
```
