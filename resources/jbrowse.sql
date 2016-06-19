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
