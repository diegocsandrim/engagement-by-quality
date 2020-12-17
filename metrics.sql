    --select uuid from  projects where kee = 'kelseyhightower:envconfig' -- "AXB_RALuGBpMHMrSj7DA"
--select analysis_uuid from  ce_activity where main_component_uuid='AXB_RALuGBpMHMrSj7DA' -- "AXB_RIlr5c8wCbutQib9

--select * from project_measures where analysis_uuid='AXB_RIlr5c8wCbutQib9'
--select * from metrics order by id

select distinct
	projects.kee,
	ce_activity.submitted_at,
	substring(substring(convert_from(ce_scanner_context.context_data, 'UTF-8'), 'sonar.analysis.contributorGain=\d*') from length('sonar.analysis.contributorGain=')+1)::integer as contributor_gain,
	metrics.name, metrics.description, metrics.short_name,
	project_measures.value, project_measures.text_value, project_measures.alert_text, project_measures.description, project_measures.variation_value_1 
from projects 
	left join ce_activity
		on projects.project_uuid=ce_activity.component_uuid
	left join ce_scanner_context
		on ce_activity.uuid = ce_scanner_context.task_uuid
	left join project_measures
		on ce_activity.analysis_uuid = project_measures.analysis_uuid
	left join metrics
		on project_measures.metric_id=metrics.id
where
	projects.kee='kelseyhightower:envconfig'
	--and ce_activity.submitted_at=1582683162603
	and metrics.name in ('complexity') -- in ('bugs',)
order by
	projects.kee,
	ce_activity.submitted_at

--metrics of interest:
-- "bugs"
-- "code_smells"
-- "cognitive_complexity"
-- "comment_lines_density"
-- "comment_lines"
-- "complexity"
-- "duplicated_blocks"
-- "duplicated_files"
-- "duplicated_lines_density"
-- "duplicated_lines"
-- "files"
-- "functions"
-- "ncloc"
-- "open_issues"
-- "sqale_debt_ratio"
-- "sqale_index"
-- "violations"
-- "vulnerabilities"


-- extract contributor gain from scanner context
select
	substring(substring(convert_from(context_data, 'UTF-8'), 'sonar.analysis.contributorGain=\d*') from length('sonar.analysis.contributorGain=')+1)::integer as contributor_gain
from
    ce_scanner_context;


task_uuid
created_at
updated_at


"AXPAWbZgdLL8_IXJ_HbQ"
1596659971718
1596659971718