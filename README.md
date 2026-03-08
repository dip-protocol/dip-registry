\# DIP Registry



Reference \*\*transparency registry\*\* for the Decision Integrity Protocol.



The registry maintains an \*\*append-only log of artifact hashes\*\* and produces \*\*Merkle inclusion proofs\*\*.



---



\# Registry Responsibilities



\* Store artifact hashes

\* Maintain append-only log

\* Compute Merkle tree

\* Produce inclusion proofs



---



\# Registry Architecture



```

artifact append

&nbsp;     ↓

append log entry

&nbsp;     ↓

update Merkle tree

&nbsp;     ↓

publish Merkle root

```



Artifacts are stored alongside the registry log.



---



\# Commands



Append artifact:



```

registry append artifact.json

```



Verify registry:



```

registry verify

```



Generate inclusion proof:



```

registry proof <artifact\_id>

```



---



\# Registry Files



```

artifacts/

log.json

merkle-root.txt

```



---



\# Proof Format



Example proof:



```json

{

&nbsp; "artifact\_hash": "...",

&nbsp; "proof\_path": \[

&nbsp;   {"hash":"...","position":"left"},

&nbsp;   {"hash":"...","position":"right"}

&nbsp; ],

&nbsp; "root": "..."

}

```



---



\# License



Apache License 2.0



