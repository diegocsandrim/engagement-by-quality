
	SELECT
		keys[1] as project_name,
		keys[2]::date as date,
		keys[3]::integer as contributor_gain,
		blocker_violations,
		bugs,
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
		effort_to_reach_maintainability_rating_a,
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
		reliability_remediation_effort,
		security_rating,
		security_remediation_effort,
		sqale_debt_ratio,
		sqale_index,
		sqale_rating,
		statements,
		violations,
		vulnerabilities
	FROM crosstab('select
	array [
	projects.kee::text,
	to_timestamp(snapshots.created_at/1000)::date::text,
	substring(substring(convert_from(ce_scanner_context.context_data, ''UTF-8''), ''sonar.analysis.contributorGain=\d*'') from length(''sonar.analysis.contributorGain='')+1)::text
	] keys,
	metrics.name as metric_name,
	replace(project_measures.value::text, ''.'', '','') AS me	tric_value
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
		''blocker_violations'',
		''bugs'',
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
		''effort_to_reach_maintainability_rating_a'',
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
		''reliability_remediation_effort'',
		''security_rating'',
		''security_remediation_effort'',
		''sqale_debt_ratio'',
		''sqale_index'',
		''sqale_rating'',
		''statements'',
		''violations'',
		''vulnerabilities''
		)
	and projects.scope=''PRJ''
order by
	projects.kee,
	snapshots.created_at,
	metrics.name') 
		AS final_result(keys TEXT[], blocker_violations TEXT,
	bugs TEXT,
	classes TEXT,
	code_smells TEXT,
	cognitive_complexity TEXT,
	comment_lines TEXT,
	comment_lines_density TEXT,
	complexity TEXT,
	critical_violations TEXT,
	duplicated_blocks TEXT,
	duplicated_files TEXT,
	duplicated_lines TEXT,
	duplicated_lines_density TEXT,
	effort_to_reach_maintainability_rating_a TEXT,
	file_complexity TEXT,
	files TEXT,
	functions TEXT,
	info_violations TEXT,
	lines TEXT,
	major_violations TEXT,
	minor_violations TEXT,
	ncloc TEXT,
	open_issues TEXT,
	reliability_rating TEXT,
	reliability_remediation_effort TEXT,
	security_rating TEXT,
	security_remediation_effort TEXT,
	sqale_debt_ratio TEXT,
	sqale_index TEXT,
	sqale_rating TEXT,
	statements TEXT,
	violations TEXT,
	vulnerabilities TEXT);