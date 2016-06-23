-- name: total-feature-count
SELECT COUNT(feat.feature_id) featcount
FROM feature feat
JOIN featureloc floc
ON feat.feature_id = floc.feature_id
JOIN feature srcfeat
ON floc.srcfeature_id = srcfeat.feature_id
JOIN organism org
ON srcfeat.organism_id = org.organism_id
JOIN cvterm cvt
ON srcfeat.type_id = cvt.cvterm_id
JOIN cv 
ON cvt.cv_id = cv.cv_id
WHERE 
org.genus = $1
AND
org.species = $2
AND
cvt.name = $3
AND
cv.name = 'sequence'

--name: total-feature-length
WITH srcfeat AS (
  SELECT srcfeat.seqlen
    FROM feature srcfeat
  JOIN organism org
    ON srcfeat.organism_id = org.organism_id
  JOIN cvterm cvt
    ON srcfeat.type_id = cvt.cvterm_id
  JOIN cv 
    ON cvt.cv_id = cv.cv_id
  WHERE 
    org.genus = $1
    AND
    org.species = $2
    AND
    cvt.name = $3
    AND
    cv.name = 'sequence'
)
SELECT SUM(srcfeat.seqlen) featlen FROM srcfeat

-- name: total-feature-count-by-coordinates
SELECT COUNT(feat.feature_id) featcount
FROM feature feat
JOIN featureloc floc
ON feat.feature_id = floc.feature_id
JOIN feature srcfeat
ON floc.srcfeature_id = srcfeat.feature_id
WHERE 
srcfeat.uniquename = $1
AND
floc.fmin >= $2
AND
floc.fmax <= $3

-- name: overlapped-features-by-coordinates
SELECT 
    floc.fmin,
    floc.fmax,
    floc.strand,
    feat.name,
    feat.uniquename,
FROM feature srcfeat
JOIN featureloc floc
ON srcfeat.feature_id = floc.srcfeature_id
JOIN feature feat
ON floc.feature_id = feat.feature_id
JOIN cvterm
ON feat.type_id = cvterm.cvterm_id
JOIN cv
ON cvterm.cv_id = cv.cv_id
WHERE
 srcfeat.uniquename = $1
 AND
 cvterm.name = $2
 AND
 cv.name = 'sequence'
 AND
(
    (
        floc.fmin >= $3
        AND
        floc.fmin <= $4
        AND
        floc.fmax > $4
    )
    OR
    (
        floc.fmin < $3
        AND
        floc.fmax >= $3
        AND
        floc.fmax <= $4
    )
    OR
    (
        floc.fmin >= $3
        AND
        floc.fmax <= $4
    )
)
                         

-- name: get-sub-sequence
SELECT SUBSTR(feat.residues, $2, $3)
FROM feature feat
WHERE 
feat.uniquename = $1
AND
feat.residues IS NOT NULL

-- name: get-reference-sotype 
SELECT cvterm.name
FROM feature feat
JOIN featureloc floc
ON feat.feature_id = floc.feature_id
JOIN cvterm
ON feat.type_id = cvterm.cvterm_id
JOIN cv 
ON cvterm.cv_id = cv.cv_id
JOIN organism 
ON feat.organism_id = organism.organism_id
WHERE
floc.srcfeature_id IS NULL
AND
cv.name = 'sequence'
AND
organism.genus = $1
AND
organism.species = $2
LIMIT 1

-- name: get-reference-features
SELECT 
    feat.uniquename,
    feat.name,
    floc.fmin,
    floc.fmax
FROM feature feat
JOIN featureloc floc
ON feat.feature_id = floc.feature_id
JOIN cvterm
ON feat.type_id = cvterm.cvterm_id
JOIN cv
ON cvterm.cv_id = cv.cv_id
JOIN organism 
ON feat.organism_id = organism.organism_id
WHERE
organism.genus = $1
AND
organism.species = $2
AND
cvterm.name = $3
AND
cv.name = 'sequence'

-- name: exact-match-name
SELECT
FROM feature feat
JOIN featureloc floc
ON feat.feature_id = floc.feature_id
JOIN feature srcfeat
ON floc.srcfeature_id = srcfeat.feature_id
JOIN organism 
ON feat.organism_id = organism.organism_id
organism.genus = $1
AND
organism.species = $2
AND
feat.name = $3
AND
floc.srcfeature_id IS NOT NULL


-- name: exact-match-name
SELECT
FROM feature feat
JOIN featureloc floc
ON feat.feature_id = floc.feature_id
JOIN feature srcfeat
ON floc.srcfeature_id = srcfeat.feature_id
JOIN organism 
ON feat.organism_id = organism.organism_id
organism.genus = $1
AND
organism.species = $2
AND
feat.name LIKE $3
AND
floc.srcfeature_id IS NOT NULL

-- name: get-organism-with-features
WITH ogroup AS (
  SELECT COUNT(feat.feature_id) featcount, feat.organism_id
    FROM feature feat
  GROUP BY feat.organism_id
  )
SELECT organism.genus, organism.species, organism.organism_id
FROM organism
JOIN featgroup
ON organism.organism_id = featgroup.organism_id

-- name: get-jbrowse-dataset-ids
SELECT jsonb_object_keys( 
        jsonb_extract_path(
            jbrowse.configuration,"general"
        ) 
    ) id
FROM jbrowse 
WHERE jbrowse.name = $1

-- name: get-each-jbrowse-dataset
SELECT jsonb_extract_path_text(
        jbrowse.configuration, "general", $2
    ), jbrowse.jbrowse_id
FROM jbrowse
WHERE jbrowse.name = $1

-- name: insert-jbrowse-organism
INSERT INTO jbrowse_organism(organism_id, jbrowse_id, dataset)
VALUES ($1, $2, $3) RETURNING jbrowse_organism_id

-- name: insert-jbrowse-track
INSERT INTO jbrowse_track(configuration, jbrowse_organism_id)
VALUES ($1, $2);

-- name: insert-jbrowse-track-with-type
INSERT INTO jbrowse_track(configuration, type_id, jbrowse_organism_id)
VALUES ($1, $2, $3);

-- name: feature-exists
SELECT COUNT(feat.feature_id) FROM feat
JOIN cvterm
ON feat.type_id = cvterm.cvterm_id
JOIN cv
ON cvterm.cv_id = cv.cv_id
WHERE
feat.organism_id = $1
AND
cvterm.name = $2
AND 
cv.name = 'sequence'

-- name: feature-with-subfeat-exist
SELECT COUNT(feat.feature_id) FROM feat
JOIN cvterm ftype
ON feat.type_id = ftype.cvterm_id
JOIN cv scv
ON ftype.cv_id = scv.cv_id
JOIN feature_relationship frel
ON feat.feature_id = frel.object_id
JOIN feature subfeat
ON frel.subject_id = subfeat.feature_id
JOIN cvterm rtype
ON frel.type_id = rtype.cvterm_id
JOIN cv rcv
ON rtype.cv_id = rcv.cv_id
JOIN cvterm subftype
ON subfeat.type_id = subftype.cvterm_id
JOIN cv subfcv
ON subftype.cv_id = subfcv.cv_id
WHERE
feat.organism_id = $1
AND
ftype.name = $2
AND
subftype.name = $3
AND 
scv.name = 'sequence'
AND
rcv.name = 'ro'
AND
subfcv.name = 'sequence'
