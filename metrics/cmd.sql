COPY (
	SELECT
		keys[1] as project_name,
		keys[2]::date as date,
		keys[3]::integer as contributor_gain,
		classes,
		code_smells,
		cognitive_complexity,
		comment_lines,
		comment_lines_density,
		complexity,
		critical_violations,
		duplicated_blocks,
		duplicated_files,
		duplicated_lines,
		duplicated_lines_density,
		file_complexity,
		files,
		functions,
		info_violations,
		lines,
		major_violations,
		minor_violations,
		ncloc,
		open_issues,
		reliability_rating,
		security_rating,
		sqale_debt_ratio,
		sqale_index,
		sqale_rating,
		statements,
		violations
	FROM crosstab('select
	array [
	projects.kee::text,
	to_timestamp(snapshots.created_at/1000)::date::text,
	substring(substring(convert_from(ce_scanner_context.context_data, ''UTF-8''), ''sonar.analysis.contributorGain=\d*'') from length(''sonar.analysis.contributorGain='')+1)::text
	] keys,
	metrics.name as metric_name,
	project_measures.value AS metric_value
from projects 
	left join ce_activity
		on projects.project_uuid=ce_activity.component_uuid
	left join snapshots
		on snapshots.uuid = ce_activity.analysis_uuid
	left join ce_scanner_context
		on ce_activity.uuid = ce_scanner_context.task_uuid
	left join project_measures
		on ce_activity.analysis_uuid = project_measures.analysis_uuid
	left join metrics
		on project_measures.metric_id=metrics.id
where
	metrics.name in (
		''classes'',
		''code_smells'',
		''cognitive_complexity'',
		''comment_lines'',
		''comment_lines_density'',
		''complexity'',
		''critical_violations'',
		''duplicated_blocks'',
		''duplicated_files'',
		''duplicated_lines'',
		''duplicated_lines_density'',
		''file_complexity'',
		''files'',
		''functions'',
		''info_violations'',
		''lines'',
		''major_violations'',
		''minor_violations'',
		''ncloc'',
		''open_issues'',
		''reliability_rating'',
		''security_rating'',
		''sqale_debt_ratio'',
		''sqale_index'',
		''sqale_rating'',
		''statements'',
		''violations''
		)
	and projects.scope=''PRJ''
order by
	projects.kee,
	snapshots.created_at,
	metrics.name') 
		AS final_result(keys TEXT[],
	classes NUMERIC,
	code_smells NUMERIC,
	cognitive_complexity NUMERIC,
	comment_lines NUMERIC,
	comment_lines_density NUMERIC,
	complexity NUMERIC,
	critical_violations NUMERIC,
	duplicated_blocks NUMERIC,
	duplicated_files NUMERIC,
	duplicated_lines NUMERIC,
	duplicated_lines_density NUMERIC,
	file_complexity NUMERIC,
	files NUMERIC,
	functions NUMERIC,
	info_violations NUMERIC,
	lines NUMERIC,
	major_violations NUMERIC,
	minor_violations NUMERIC,
	ncloc NUMERIC,
	open_issues NUMERIC,
	reliability_rating NUMERIC,
	security_rating NUMERIC,
	sqale_debt_ratio NUMERIC,
	sqale_index NUMERIC,
	sqale_rating NUMERIC,
	statements NUMERIC,
	violations NUMERIC)
)
TO '/tmp/out.csv' (format csv, delimiter ';', HEADER);
